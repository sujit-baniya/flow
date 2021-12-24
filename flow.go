package flow

import (
	"fmt"
	"github.com/pkg/errors"
)

type Flow struct {
	Key       string `json:"key"`
	firstNode Node
	Branch    bool `json:"branch"`
	nodes     map[string]Node
	inVertex  map[string]bool
	outVertex map[string]bool
	Error     error
}

func New() *Flow {
	return &Flow{
		nodes:     make(map[string]Node),
		inVertex:  make(map[string]bool),
		outVertex: make(map[string]bool),
	}
}

func (f *Flow) Node(vertex string, handler Handler) *Flow {
	f.nodes[vertex] = &Vertex{Key: vertex, handler: handler, edges: make(map[string]Node), branches: make(map[string]Node), ConditionalNodes: make(map[string]string)}
	return f
}

func (f *Flow) Edge(inVertex, outVertex string) *Flow {
	var outNode, inNode Node
	var okOutNode, okInNode bool
	outNode, okOutNode = f.nodes[outVertex]
	inNode, okInNode = f.nodes[inVertex]
	if !okOutNode {
		f.Error = errors.New(fmt.Sprintf("Output Vertex with key %s doesn't exist", outVertex))
		return f
	}
	if !okInNode {
		f.Error = errors.New(fmt.Sprintf("Input Vertex with key %s doesn't exist", inVertex))
		return f
	}
	f.inVertex[inVertex] = true
	f.outVertex[outVertex] = true
	inOk := f.inVertex[inVertex]
	outOk := f.outVertex[inVertex]
	if inOk && !outOk {
		f.firstNode = f.nodes[inVertex]
	}
	if okInNode && okOutNode {
		inNode.AddEdge(outNode)
	}
	return f
}

func (f *Flow) ConditionalNode(vertex string, handler Handler, conditions map[string]string) *Flow {
	branches := make(map[string]Node)
	node := &Vertex{
		Key:              vertex,
		Branch:           true,
		handler:          handler,
		ConditionalNodes: conditions,
	}
	for condition, nodeKey := range conditions {
		f.outVertex[nodeKey] = true
		if n, ok := f.nodes[nodeKey]; ok {
			branches[condition] = n
		}
	}
	node.branches = branches
	f.nodes[vertex] = node
	return f
}

func (f *Flow) Process(data DataSource) (DataSource, error) {
	if f.Error != nil {
		return DataSource{}, f.Error
	}
	if f.firstNode == nil {
		return DataSource{}, errors.New("No edges defined")
	}
	return f.firstNode.Process(data)
}

func (f *Flow) IsBranch() bool {
	return f.Branch
}

func (f *Flow) GetKey() string {
	return f.Key
}

func (f *Flow) AddEdge(node Node) {
	f.nodes[node.GetKey()] = node
}

func (f *Flow) SubFlow(flow *Flow) *Flow {
	f.nodes[flow.GetKey()] = flow
	return f
}