package flow

import (
	"context"
	"encoding/json"
)

type Loop struct {
	Key       string `json:"key"`
	SingleKey string `json:"child_key"`
	handler   Handler
	Single    Node
}

func (l *Loop) Process(ctx context.Context, data Data) (Data, error) {
	response, err := l.handler(ctx, data)
	if err != nil {
		return data, err
	}
	var rs []interface{}
	err = json.Unmarshal(response.Payload, &rs)
	for _, single := range rs {
		payload, _ := json.Marshal(single)
		dataPayload := data
		dataPayload.Payload = payload
		_, err = l.Single.Process(ctx, dataPayload)
		if err != nil {
			return data, err
		}
	}
	return response, err
}

func (l *Loop) AddEdge(node Node) {

}
func (l *Loop) IsBranch() bool {
	return false
}

func (l *Loop) GetKey() string {
	return l.Key
}

func (l *Loop) GetSingleKey() string {
	return l.SingleKey
}
