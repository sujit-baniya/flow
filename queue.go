package flow

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrQueueShuttingDown = errors.New("queue is shutting down; new tasks are not being accepted")

type Queue struct {
	mutex sync.Mutex
	name  string
	next  time.Time
	now   func() time.Time
	tasks []*Task
	wg    sync.WaitGroup

	accept   int32
	shutdown chan struct{}
	started  chan struct{}
	wake     chan struct{}
}

// NewQueue Creates a new task queue. The name of the task queue is used in Prometheus
// label names and must match [a-zA-Z0-9:_] (snake case is used by convention).
func NewQueue(name string) *Queue {
	return &Queue{
		name: name,
		now: func() time.Time {
			return time.Now().UTC()
		},
		accept:   1,
		shutdown: make(chan struct{}),
		started:  make(chan struct{}, 1),
	}
}

// Now Sets the function the queue will use to obtain the current time.
func (q *Queue) Now(now func() time.Time) {
	q.now = now
}

// Enqueue Enqueues a task.
//
// An error will only be returned if the queue has been shut down.
func (q *Queue) Enqueue(t *Task) error {
	if atomic.LoadInt32(&q.accept) == 0 {
		return ErrQueueShuttingDown
	}

	q.mutex.Lock()
	q.tasks = append(q.tasks, t)
	if q.wake != nil {
		// Runs asynchronously to avoid deadlocking if a task submits another task
		go func() {
			q.wake <- struct{}{}
		}()
	}
	q.mutex.Unlock()
	return nil
}

// Submit Creates and enqueues a new task, returning the new task. Note that the
// caller cannot customize settings on the task without creating a race
// condition; so attempting to will panic. See NewTask and (*Queue).Enqueue to
// create tasks with customized options.
//
// An error will only be returned if the queue has been shut down.
func (q *Queue) Submit(fn TaskFunc) (*Task, error) {
	t := NewTask(fn)
	t.immutable = true
	err := q.Enqueue(t)
	return t, err
}

// Dispatch Attempts any tasks which are due and updates the task schedule. Returns true
// if there is more work to do.
func (q *Queue) Dispatch(ctx context.Context) bool {
	next := time.Unix(1<<63-62135596801, 999999999) // "max" time
	now := q.now()

	// In order to avoid deadlocking if a task queues another task, we make a
	// copy of the task list and release the mutex while executing them.
	q.mutex.Lock()
	tasks := make([]*Task, len(q.tasks))
	copy(tasks, q.tasks)
	q.mutex.Unlock()

	for _, task := range tasks {
		due := task.NextAttempt().Before(now)
		if due {
			_, _ = task.Attempt(ctx)
		}
		if !task.Done() && task.NextAttempt().Before(next) {
			next = task.NextAttempt()
		}
	}

	q.mutex.Lock()
	newTasks := make([]*Task, 0, len(q.tasks))
	for _, task := range q.tasks {
		if !task.Done() {
			newTasks = append(newTasks, task)
		}
	}
	q.tasks = newTasks
	q.mutex.Unlock()

	q.next = next
	return len(newTasks) != 0
}

func (q *Queue) run(ctx context.Context) {
	q.mutex.Lock()
	if q.wake != nil {
		panic(errors.New("this queue is already running on another goroutine"))
	}

	q.wake = make(chan struct{})
	q.mutex.Unlock()

	for {
		more := q.Dispatch(ctx)
		if atomic.LoadInt32(&q.accept) == 0 && !more {
			return
		}

		select {
		case <-time.After(q.next.Sub(q.now())):
			break
		case <-ctx.Done():
			return
		case <-q.wake:
			break
		case <-q.shutdown:
			atomic.StoreInt32(&q.accept, 0)
			break
		}
	}
}

// Run the task queue. Blocks until the context is cancelled.
func (q *Queue) Run(ctx context.Context) {
	select {
	case <-q.started:
		panic(errors.New("this queue is already started on another goroutine"))
	default:
		q.run(ctx)
	}
}

// Start the task queue in the background. If you wish to use the warm
// shutdown feature, you must use Start, not Run.
func (q *Queue) Start(ctx context.Context) {
	q.wg.Add(1)

	select {
	case q.started <- struct{}{}:
		go func() {
			q.run(ctx)
			q.wg.Done()
		}()
	default:
		panic(errors.New("this queue is already started on another goroutine"))
	}
}

// Shutdown Stops accepting new tasks and blocks until all already-queued tasks are
// complete. The queue must have been started with Start, not Run.
func (q *Queue) Shutdown() {
	select {
	case <-q.started:
		q.shutdown <- struct{}{}
		q.wg.Wait()
	default:
		panic(errors.New("attempted warm shutdown on queue which was not run with queue.Start(ctx)"))
	}
}

// Join Shuts down any number of work queues in parallel and blocks until they're
// all finished.
func Join(queues ...*Queue) {
	var wg sync.WaitGroup
	wg.Add(len(queues))
	for _, q := range queues {
		go func(q *Queue) {
			q.Shutdown()
			wg.Done()
		}(q)
	}
	wg.Wait()
}
