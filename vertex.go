package flow

import (
	"context"
	"encoding/json"
	"errors"
)

type Vertex struct {
	Key              string `json:"key"`
	Type             string `json:"type"`
	handler          Handler
	edges            map[string]Node
	branches         map[string]Node
	loops            map[string]Node
	ConditionalNodes map[string]string `json:"conditional_nodes"`
}

func (v *Vertex) Process(ctx context.Context, data Data) (Data, error) {
	if v.GetType() == "Branch" && len(v.ConditionalNodes) == 0 {
		return Data{}, errors.New("required at least one condition for branch")
	}
	if data.visitedVertices == nil {
		data.visitedVertices = make(map[string]int)
	}
	data.CurrentVertex = v.GetKey()
	data.visitedVertices[v.GetKey()]++
	response, err := v.handler(ctx, data)
	if err != nil {
		return Data{}, err
	}
	if v.Type == "Loop" {
		var rs []interface{}
		err = json.Unmarshal(response.Payload, &rs)
		for _, single := range rs {
			payload, _ := json.Marshal(single)
			dataPayload := data
			dataPayload.Payload = payload
			for _, loop := range v.loops {
				response, err = loop.Process(ctx, dataPayload)
				if err != nil {
					return Data{}, err
				}
			}
		}
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

func (v *Vertex) GetType() string {
	return v.Type
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
