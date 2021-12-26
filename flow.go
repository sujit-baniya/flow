package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type Flow struct {
	Key       string `json:"key"`
	Branch    bool   `json:"branch"`
	Error     error
	firstNode Node
	nodes     map[string]Node
	inVertex  map[string]bool
	outVertex map[string]bool
	raw       RawFlow
}

type RawFlow struct {
	Nodes    []string   `json:"nodes"`
	Loops    [][]string `json:"Loops"`
	Branches []Branch   `json:"branches"`
	Edges    [][]string `json:"edges"`
}

type Branch struct {
	Key              string            `json:"key"`
	ConditionalNodes map[string]string `json:"conditional_nodes"`
}

func New(raw ...Payload) *Flow {
	f := &Flow{
		nodes:     make(map[string]Node),
		inVertex:  make(map[string]bool),
		outVertex: make(map[string]bool),
	}
	if len(raw) == 0 {
		return f
	}
	var rawFlow RawFlow
	err := json.Unmarshal(raw[0], &rawFlow)
	if err != nil {
		f.Error = err
		return f
	}
	f.raw = rawFlow
	return f
}

func (f *Flow) Node(vertex string) *Flow {
	f.raw.Nodes = append(f.raw.Nodes, vertex)
	return f
}

func (f *Flow) Edge(inVertex, outVertex string) *Flow {
	f.raw.Edges = append(f.raw.Edges, []string{inVertex, outVertex})
	return f
}

func (f *Flow) ConditionalNode(vertex string, conditions map[string]string) *Flow {
	branch := Branch{
		Key:              vertex,
		ConditionalNodes: conditions,
	}
	f.raw.Branches = append(f.raw.Branches, branch)
	return f
}

func (f *Flow) Loop(inVertex, childVertex string) *Flow {
	f.raw.Loops = append(f.raw.Loops, []string{inVertex, childVertex})
	return f
}

func (f *Flow) loop(inVertex, childVertex string, inHandler Handler) *Flow {
	f.outVertex[childVertex] = true
	loop := &Loop{
		Key:       inVertex,
		SingleKey: childVertex,
		handler:   inHandler,
		Single:    f.nodes[childVertex],
	}
	f.nodes[inVertex] = loop
	return f
}

func (f *Flow) Process(ctx context.Context, data Data) (Data, error) {
	if f.Error != nil {
		return Data{}, f.Error
	}
	if f.firstNode == nil {
		return Data{}, errors.New("No edges defined")
	}
	return f.firstNode.Process(ctx, data)
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

func (f *Flow) Build() *Flow {
	var noNodes, noEdges bool
	for _, node := range f.raw.Nodes {
		handler := GetNodeHandler(node)
		if handler != nil {
			f.node(node, handler)
		} else {
			f.Error = errors.New(fmt.Sprintf("No node handler defined for key '%s'", node))
			return f
		}
	}
	for _, branch := range f.raw.Branches {
		branchHandler := GetBranchHandler(branch.Key)
		if branchHandler != nil {
			for _, node := range branch.ConditionalNodes {
				nodeHandler := GetNodeHandler(node)
				if nodeHandler != nil {
					f.node(node, nodeHandler)
				}
			}
		}
	}
	if len(f.raw.Edges) == 0 {
		noEdges = true
	}
	for _, edge := range f.raw.Edges {
		inVertex := edge[0]
		outVertex := edge[1]
		inNodeHandler := GetNodeHandler(inVertex)
		if inNodeHandler != nil {
			f.node(inVertex, inNodeHandler)
		}
		outNodeHandler := GetNodeHandler(outVertex)
		if outNodeHandler != nil {
			f.node(outVertex, outNodeHandler)
		}
	}

	for _, loop := range f.raw.Loops {
		vertex := loop[0]
		childVertex := loop[1]
		loopHandler := GetNodeHandler(vertex)
		if loopHandler != nil {
			f.loop(vertex, childVertex, loopHandler)
		}
	}
	for _, branch := range f.raw.Branches {
		branchHandler := GetBranchHandler(branch.Key)
		if branchHandler == nil {
			f.Error = errors.New(fmt.Sprintf("No branch handler defined for key '%s'", branch.Key))
			return f
		}
		f.conditionalNode(branch.Key, branchHandler, branch.ConditionalNodes)
	}

	for _, edge := range f.raw.Edges {
		inVertex := edge[0]
		outVertex := edge[1]
		f.edge(inVertex, outVertex)
	}
	if noEdges || noNodes {
		f.Error = errors.New("No vertex or edges are defined")
	}
	return f
}

func (f *Flow) conditionalNode(vertex string, handler Handler, conditions map[string]string) *Flow {
	if _, ok := f.nodes[vertex]; !ok {
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
	}
	return f
}

func (f *Flow) node(vertex string, handler Handler) *Flow {
	if _, ok := f.nodes[vertex]; !ok {
		f.nodes[vertex] = &Vertex{
			Key:              vertex,
			ConditionalNodes: make(map[string]string),
			handler:          handler,
			edges:            make(map[string]Node),
			branches:         make(map[string]Node),
		}
	}

	return f
}

func (f *Flow) edge(inVertex, outVertex string) *Flow {
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

var BranchList = map[string]Handler{}

func AddBranch(node string, handler Handler) {
	BranchList[node] = handler
}

func GetBranchHandler(node string) Handler {
	return BranchList[node]
}
