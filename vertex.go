package flow

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/sync/errgroup"
	"reflect"
)

type Vertex struct {
	Key              string            `json:"key"`
	Type             string            `json:"type"`
	ConditionalNodes map[string]string `json:"conditional_nodes"`
	handler          Handler
	edges            map[string]Node
	branches         map[string]Node
	loops            map[string]Node
}

func loop(ctx context.Context, loops map[string]Node, data Data, response Data) ([]interface{}, error) {
	g, ctx := errgroup.WithContext(ctx)
	var rs, results []interface{}
	err := json.Unmarshal(response.Payload, &rs)
	if err != nil {
		return nil, err
	}
	for _, single := range rs {
		single := single
		g.Go(func() error {
			payload, err := json.Marshal(single)
			if err != nil {
				return err
			}
			dataPayload := data
			dataPayload.Payload = payload
			for _, loop := range loops {
				resp, err := loop.Process(ctx, dataPayload)
				if err != nil {
					return err
				}
				err = json.Unmarshal(resp.Payload, &single)
				if err != nil {
					return err
				}
				results = append(results, single)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func (v *Vertex) Process(ctx context.Context, data Data) (Data, error) {
	if v.GetType() == "Branch" && len(v.ConditionalNodes) == 0 {
		return Data{}, errors.New("required at least one condition for branch")
	}
	data.CurrentVertex = v.GetKey()
	response, err := v.handler(ctx, data)
	if err != nil {
		return Data{}, err
	}
	if v.Type == "Loop" {
		result, err := loop(ctx, v.loops, data, response)
		if err != nil {
			return Data{}, err
		}
		/*var rs []interface{}
		var result []interface{}
		err = json.Unmarshal(response.Payload, &rs)
		for _, single := range rs {
			payload, _ := json.Marshal(single)
			dataPayload := data
			dataPayload.Payload = payload
			for _, loop := range v.loops {
				resp, err := loop.Process(ctx, dataPayload)
				if err != nil {
					return Data{}, err
				}
				err = json.Unmarshal(resp.Payload, &single)
				if err != nil {
					return Data{}, err
				}
				result = append(result, single)
			}
		}*/
		tmp, err := json.Marshal(result)
		if err != nil {
			return Data{}, err
		}
		response.Payload = tmp
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

func clone(data interface{}) interface{} {
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(data).Elem().Type()).Interface()
	}
	return reflect.New(reflect.TypeOf(data)).Elem().Interface()
}
