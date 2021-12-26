package flow

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type Vertex struct {
	Key              string `json:"key"`
	Branch           bool   `json:"branch"`
	handler          Handler
	edges            map[string]Node
	branches         map[string]Node
	loops            map[string]Node
	ConditionalNodes map[string]string `json:"conditional_nodes"`
}

func (v *Vertex) Process(ctx context.Context, data Data) (Data, error) {
	if v.IsBranch() && len(v.ConditionalNodes) == 0 {
		return Data{}, errors.New("Required at least one condition for branch")
	}
	if data.visitedVertices == nil {
		data.visitedVertices = make(map[string]int)
	}
	data.CurrentVertex = v.GetKey()
	fmt.Println(data.CurrentVertex)
	data.visitedVertices[v.GetKey()]++
	response, err := v.handler(ctx, data)
	if err != nil {
		return Data{}, err
	}

	if val, ok := v.branches[response.GetStatus()]; ok {
		response, err = val.Process(ctx, response)
	}
	for _, edge := range v.edges {
		response, err = edge.Process(ctx, response)
		if err != nil {
			return Data{}, err
		}
	}
	return response, err
}

func (v *Vertex) IsBranch() bool {
	return v.Branch
}

func (v *Vertex) GetKey() string {
	return v.Key
}

func (v *Vertex) AddEdge(node Node) {
	if v.edges == nil {
		v.edges = make(map[string]Node)
	}
	v.edges[node.GetKey()] = node
}
