package flow

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type Node interface {
	Process(data DataSource) (DataSource, error)
	AddEdge(node Node)
	IsBranch() bool
	GetKey() string
}

type Payload []byte

type DataSource struct {
	RequestID     string  `json:"request_id"`
	Payload       Payload `json:"payload"`
	Status        string  `json:"status"`
	CurrentVertex string  `json:"current_vertex"`
	FailedReason  error   `json:"failed_reason"`
}

func (d DataSource) ConvertTo(rs interface{}) error {
	return json.Unmarshal(d.Payload, rs)
}

func (d DataSource) ToString() string {
	return string(d.Payload)
}

func (d DataSource) GetStatus() string {
	return d.Status
}

type Handler func(data DataSource) (DataSource, error)

type Vertex struct {
	Key              string `json:"key"`
	Branch           bool   `json:"branch"`
	handler          Handler
	edges            map[string]Node
	branches         map[string]Node
	ConditionalNodes map[string]string `json:"conditional_nodes"`
}

func (v *Vertex) Process(data DataSource) (DataSource, error) {
	if v.IsBranch() && len(v.ConditionalNodes) == 0 {
		return DataSource{}, errors.New("Required at least one condition for branch")
	}
	data.CurrentVertex = v.GetKey()
	fmt.Println(data.CurrentVertex)
	response, err := v.handler(data)
	if err != nil {
		return DataSource{}, err
	}

	if val, ok := v.branches[response.GetStatus()]; ok {
		response, err = val.Process(response)
	}
	for _, edge := range v.edges {
		response, err = edge.Process(response)
		if err != nil {
			return DataSource{}, err
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
