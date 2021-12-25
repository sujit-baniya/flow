package flow

import (
	"encoding/json"
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
	raw       RawFlow
}

type RawFlow struct {
	Nodes    []string   `json:"nodes"`
	SubFlows []string   `json:"sub_flows"`
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
	json.Unmarshal(raw[0], &rawFlow)
	f.raw = rawFlow
	return f
}

func (f *Flow) Node(vertex string) *Flow {
	f.raw.Nodes = append(f.raw.Nodes, vertex)
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

func (f *Flow) Edge(inVertex, outVertex string) *Flow {
	f.raw.Edges = append(f.raw.Edges, []string{inVertex, outVertex})
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

func (f *Flow) ConditionalNode(vertex string, conditions map[string]string) *Flow {
	branch := Branch{
		Key:              vertex,
		ConditionalNodes: conditions,
	}
	f.raw.Branches = append(f.raw.Branches, branch)
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

func (f *Flow) SubFlow(flow string) *Flow {
	f.raw.SubFlows = append(f.raw.SubFlows, flow)
	return f
}

func (f *Flow) subFlow(flow *Flow) *Flow {
	f.nodes[flow.GetKey()] = flow
	return f
}

func (f *Flow) Build() *Flow {
	var noNodes, noEdges bool
	for _, node := range f.raw.Nodes {
		handler := GetNodeHandler(node)
		if handler != nil {
			f.node(node, handler)
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
	for _, branch := range f.raw.Branches {
		branchHandler := GetBranchHandler(branch.Key)
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

var NodeList = map[string]Handler{}

func AddNode(node string, handler Handler) {
	NodeList[node] = handler
}

func GetNodeHandler(node string) Handler {
	return NodeList[node]
}

func GetNodeList() []string {
	var nodes []string
	for node, _ := range NodeList {
		nodes = append(nodes, node)
	}
	return nodes
}

var BranchList = map[string]Handler{}

func AddBranch(node string, handler Handler) {
	BranchList[node] = handler
}

func GetBranchList() []string {
	var branches []string
	for branch, _ := range BranchList {
		branches = append(branches, branch)
	}
	return branches
}

func GetBranchHandler(node string) Handler {
	return BranchList[node]
}

var FlowList = map[string]*Flow{}

func AddFlow(node string, flow *Flow) {
	FlowList[node] = flow
}

func GetFlowList() []string {
	var flows []string
	for flow, _ := range FlowList {
		flows = append(flows, flow)
	}
	return flows
}

func GetFlow(node string) *Flow {
	return FlowList[node]
}
