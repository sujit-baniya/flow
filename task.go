package flow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

var (
	// ErrAlreadyComplete Returned when a task is attempted which was already successfully completed.
	ErrAlreadyComplete = errors.New("this task was already successfully completed once")

	// ErrDoNotReattempt If this is returned from a task function, the task shall not be re-attempted.
	ErrDoNotReattempt = errors.New("this task should not be re-attempted")

	// ErrMaxRetriesExceeded This task has been attempted too many times.
	ErrMaxRetriesExceeded = errors.New("the maximum retries for this task has been exceeded")

	// Now Set this function to influence the clock that will be used for
	// scheduling re-attempts.
	Now = func() time.Time {
		return time.Now().UTC()
	}
)

type TaskFunc func(ctx context.Context) error

// Task Stores state for a task which shall be or has been executed. Each task may
// only be executed successfully once.
type Task struct {
	Metadata map[string]interface{}

	after       func(ctx context.Context, task *Task)
	attempts    int
	done        bool
	err         error
	fn          TaskFunc
	nextAttempt time.Time

	baseDuration  time.Duration
	immutable     bool
	jitter        bool
	maxAttempts   int
	maxTimeout    time.Duration
	sleepDuration time.Duration
	within        time.Duration
}

// NewTask Creates a new task for a given function.
func NewTask(fn TaskFunc) *Task {
	return &Task{
		Metadata: make(map[string]interface{}),

		baseDuration:  time.Minute,
		fn:            fn,
		jitter:        true,
		maxAttempts:   1,
		maxTimeout:    30 * time.Minute,
		sleepDuration: time.Minute,
	}
}

// Attempt to execute this task.
//
// If successful, the zero time and nil are returned.
//
// Otherwise, the error returned from the task function is returned to the
// caller. If an error is returned for which errors.Is(err, ErrDoNotReattempt)
// is true, the caller should not call Attempt again.
func (t *Task) Attempt(ctx context.Context) (time.Time, error) {
	if t.done {
		if t.err == nil {
			return time.Time{}, ErrAlreadyComplete
		}
		return time.Time{}, t.err
	}

	t.attempts += 1
	if t.attempts > t.maxAttempts {
		t.err = ErrMaxRetriesExceeded
		t.done = true
		if t.after != nil {
			t.after(ctx, t)
			t.after = nil
		}
		return time.Time{}, ErrMaxRetriesExceeded
	}

	if errors.Is(t.err, ErrDoNotReattempt) {
		return time.Time{}, t.err
	}

	if t.within != time.Duration(0) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t.within)
		defer cancel()
	}

	func() {
		defer func() {
			if err := recover(); err != nil {
				t.err = fmt.Errorf("panic: %v", err)
			}
		}()
		t.err = t.fn(ctx)
	}()

	if t.err == nil {
		t.done = true
		if t.after != nil {
			t.after(ctx, t)
			t.after = nil
		}
		return time.Time{}, nil
	}

	if t.jitter {
		rand.Seed(time.Now().UnixNano())
		max := int((t.sleepDuration*3 - t.baseDuration + time.Minute).Minutes())
		sleep := time.Duration(rand.Intn(max))*time.Minute + t.baseDuration
		if t.sleepDuration = sleep; t.sleepDuration > t.maxTimeout {
			t.sleepDuration = t.maxTimeout
		}
	} else {
		sleep := math.Pow(2, float64(t.attempts))
		if sleep > t.maxTimeout.Minutes() || sleep < 0 {
			t.sleepDuration = t.maxTimeout
		} else {
			t.sleepDuration = time.Duration(sleep) * time.Minute
		}
	}

	t.nextAttempt = Now().Add(t.sleepDuration)

	log.Printf("Attempt %d/%d failed (%v), retrying in %s",
		t.Attempts(), t.maxAttempts, t.err, t.sleepDuration.String())
	return t.nextAttempt, t.err
}

// Retries Set the maximum number of retries on failure, or -1 to attempt indefinitely.
// By default, tasks are not retried on failure.
func (t *Task) Retries(n int) *Task {
	if n < -1 {
		panic(errors.New("invalid input to Task.Retries"))
	}
	if t.immutable {
		panic(errors.New("attempted to configure immutable task"))
	}
	t.maxAttempts = n
	return t
}

// MaxTimeout Sets the maximum timeout between retries, or zero to exponentially increase
// the timeout indefinitely. Defaults to 30 minutes.
func (t *Task) MaxTimeout(d time.Duration) *Task {
	if d < 0 {
		panic(errors.New("invalid timeout provided to Task.MaxTimeout"))
	}
	if t.immutable {
		panic(errors.New("attempted to configure immutable task"))
	}
	t.maxTimeout = d
	return t
}

// After Sets a function which will be executed once the task is completed,
// successfully or not. The final result (nil or an error) is passed to the
// callee.
func (t *Task) After(fn func(ctx context.Context, task *Task)) *Task {
	if t.after != nil {
		panic(errors.New("this task already has an 'After' function assigned"))
	}
	if t.immutable {
		panic(errors.New("attempted to configure immutable task"))
	}
	t.after = fn
	return t
}

// Within Specifies an upper limit for the duration of each attempt.
func (t *Task) Within(deadline time.Duration) *Task {
	if t.immutable {
		panic(errors.New("attempted to configure immutable task"))
	}
	t.within = deadline
	return t
}

// NotBefore Specifies the earliest possible time of the first execution.
func (t *Task) NotBefore(date time.Time) *Task {
	if t.immutable {
		panic(errors.New("attempted to configure immutable task"))
	}
	t.nextAttempt = date
	return t
}

// NoJitter Specifies that randomness should not be introduced into the exponential
// backoff algorithm.
func (t *Task) NoJitter() *Task {
	t.jitter = false
	return t
}

// Result Returns the result of the task. The task must have been completed for this
// to be valid.
func (t *Task) Result() error {
	if !t.done {
		panic(errors.New("(*Task).Result() called on incomplete task"))
	}
	return t.err
}

// Attempts Returns the number of times this task has been attempted
func (t *Task) Attempts() int {
	return t.attempts
}

// NextAttempt Returns the time the next attempt is scheduled for, or the zero value if it
// has not been attempted before.
func (t *Task) NextAttempt() time.Time {
	return t.nextAttempt
}

// Done Returns true if this task was completed, successfully or not.
func (t *Task) Done() bool {
	return t.done
}
