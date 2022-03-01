package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type Flow struct {
	Key       string
	Error     error
	firstNode Node
	rawNodes  map[string]Handler
	nodes     map[string]Node
	inVertex  map[string]bool
	outVertex map[string]bool
	raw       RawFlow
}

type RawFlow struct {
	Nodes    []string   `json:"nodes,omitempty"`
	Loops    [][]string `json:"loops,omitempty"`
	ForEach  []ForEach  `json:"for_each,omitempty"`
	Branches []Branch   `json:"branches,omitempty"`
	Edges    [][]string `json:"edges,omitempty"`
}

type Branch struct {
	Key              string            `json:"key"`
	ConditionalNodes map[string]string `json:"conditional_nodes"`
}

type ForEach struct {
	InVertex    string   `json:"in_vertex"`
	ChildVertex []string `json:"child_vertex"`
}

func New(raw ...[]byte) *Flow {
	f := &Flow{
		nodes:     make(map[string]Node),
		inVertex:  make(map[string]bool),
		outVertex: make(map[string]bool),
		rawNodes:  make(map[string]Handler),
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

func NewRaw(flow RawFlow) *Flow {
	return &Flow{
		nodes:     make(map[string]Node),
		inVertex:  make(map[string]bool),
		outVertex: make(map[string]bool),
		rawNodes:  make(map[string]Handler),
		raw:       flow,
	}
}

func (f *Flow) Node(vertex string) *Flow {
	f.raw.Nodes = append(f.raw.Nodes, vertex)
	return f
}

func (f *Flow) AddNode(node string, handler Handler) *Flow {
	f.rawNodes[node] = handler
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

func (f *Flow) Loop(inVertex string, childVertex ...string) *Flow {
	v := []string{inVertex}
	v = append(v, childVertex...)
	f.raw.Loops = append(f.raw.Loops, v)
	return f
}

func (f *Flow) ForEach(inVertex string, childVertex ...string) *Flow {
	forEach := ForEach{
		InVertex:    inVertex,
		ChildVertex: childVertex,
	}
	f.raw.ForEach = append(f.raw.ForEach, forEach)
	return f
}

func (f *Flow) Process(ctx context.Context, data Data) (Data, error) {
	if f.Error != nil {
		return Data{}, f.Error
	}
	if f.firstNode == nil {
		t := f.Build()
		if t.Error != nil {
			return Data{}, t.Error
		}
	}
	if f.firstNode == nil {
		return Data{}, errors.New("no edges defined")
	}
	return f.firstNode.Process(ctx, data)
}

func (f *Flow) GetType() string {
	return "Flow"
}

func (f *Flow) GetKey() string {
	return f.Key
}

func (f *Flow) AddEdge(node Node) {
	f.nodes[node.GetKey()] = node
}

func (f *Flow) WithRaw(raw RawFlow) *Flow {
	f.raw = raw
	return f
}

func (f *Flow) GetNodeHandler(node string) Handler {
	return f.rawNodes[node]
}

func (f *Flow) Build() *Flow {
	var noNodes, noEdges bool
	for _, node := range f.raw.Nodes {
		f.addNode(node)
	}
	for _, branch := range f.raw.Branches {
		branchHandler := f.GetNodeHandler(branch.Key)
		if branchHandler != nil {
			for _, node := range branch.ConditionalNodes {
				f.addNode(node)
			}
		}
	}
	if len(f.raw.Edges) == 0 {
		noEdges = true
	}
	for _, edge := range f.raw.Edges {
		f.addNode(edge[0])
		f.addNode(edge[1])
	}

	for _, loop := range f.raw.Loops {
		loopHandler := f.GetNodeHandler(loop[0])
		childVertexes := loop[1:]
		for _, v := range childVertexes {
			f.addNode(v)
		}
		if loopHandler != nil {
			f.loop(loop[0], loopHandler, childVertexes...)
		}
	}
	for _, branch := range f.raw.Branches {
		branchHandler := f.GetNodeHandler(branch.Key)
		if branchHandler == nil {
			f.Error = errors.New(fmt.Sprintf("No branch handler defined for key '%s'", branch.Key))
			return f
		}
		f.conditionalNode(branch.Key, branchHandler, branch.ConditionalNodes)
	}

	for _, edge := range f.raw.Edges {
		f.edge(edge[0], edge[1])
	}
	if noEdges || noNodes {
		f.Error = errors.New("no vertex or edges are defined")
	}
	Add(f.Key, f)
	return f
}

func (f *Flow) addNode(node string) {
	handler := f.GetNodeHandler(node)
	if handler != nil {
		f.node(node, handler)
	}
}

func (f *Flow) conditionalNode(vertex string, handler Handler, conditions map[string]string) *Flow {
	branches := make(map[string]Node)
	node := &Vertex{
		Key:              vertex,
		Type:             "Branch",
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

func (f *Flow) node(vertex string, handler Handler) *Flow {
	if _, ok := f.nodes[vertex]; !ok {
		f.nodes[vertex] = &Vertex{
			Key:              vertex,
			Type:             "Vertex",
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

func (f *Flow) loop(inVertex string, inHandler Handler, childVertex ...string) *Flow {
	childVertexes := make(map[string]Node)
	for _, v := range childVertex {
		f.outVertex[v] = true
		childVertexes[v] = f.nodes[v]
	}

	loop := &Vertex{
		Key:     inVertex,
		Type:    "Loop",
		loops:   childVertexes,
		handler: inHandler,
	}
	f.nodes[inVertex] = loop
	return f
}

var flowList = map[string]*Flow{}

func Add(key string, flow *Flow) {
	flowList[key] = flow
}

func Get(key string) *Flow {
	return flowList[key]
}

func All() map[string]*Flow {
	return flowList
}
