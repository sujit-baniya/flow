package flow

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"golang.org/x/sync/errgroup"
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

func merge(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
	for k, m := range map2 {
		if _, ok := map1[k]; !ok {
			map1[k] = m
		}
	}
	return map1
}

func (v *Vertex) loop(ctx context.Context, loops map[string]Node, data Data, response Data) ([]interface{}, error) {
	g, ctx := errgroup.WithContext(ctx)
	result := make(chan interface{})
	var rs, results []interface{}
	err := json.Unmarshal(response.Payload, &rs)
	if err != nil {
		return nil, err
	}
	for _, single := range rs {
		single := single
		g.Go(func() error {
			var payload []byte
			currentData := make(map[string]interface{})
			switch s := single.(type) {
			case map[string]interface{}:
				currentData = s
			}
			if currentData != nil {
				payload, err = json.Marshal(currentData)
				if err != nil {
					return err
				}
			} else {
				payload, err = json.Marshal(single)
				if err != nil {
					return err
				}
			}
			dataPayload := data
			dataPayload.Payload = payload
			var responseData map[string]interface{}
			for _, loop := range loops {
				resp, err := loop.Process(ctx, dataPayload)
				resp.FailedReason = err
				if err != nil {
					return err
				}
				err = json.Unmarshal(resp.Payload, &responseData)
				if err != nil {
					return err
				}
				currentData = merge(currentData, responseData)
			}
			payload, err = json.Marshal(currentData)
			if err != nil {
				return err
			}
			dataPayload.Payload = payload
			err = json.Unmarshal(dataPayload.Payload, &single)
			if err != nil {
				return err
			}
			select {
			case result <- single:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}
	go func() {
		g.Wait()
		close(result)
	}()
	for ch := range result {
		results = append(results, ch)
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func (v *Vertex) Process(ctx context.Context, data Data) (Data, error) {
	if v.GetType() == "Branch" && len(v.ConditionalNodes) == 0 {
		return data, errors.New("required at least one condition for branch")
	}
	response, err := v.handler(ctx, data)
	if err != nil {
		return data, err
	}
	if v.Type == "Loop" {
		result, err := v.loop(ctx, v.loops, data, response)
		if err != nil {
			return data, err
		}
		tmp, err := json.Marshal(result)
		if err != nil {
			return data, err
		}
		response.Payload = tmp
	}
	if val, ok := v.branches[response.GetStatus()]; ok {
		response, err = val.Process(ctx, response)
		response.FailedReason = err
	}
	for _, edge := range v.edges {
		response, err = edge.Process(ctx, response)
		response.FailedReason = err
		if err != nil {
			return data, err
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
