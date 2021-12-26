package main

import (
	"context"
	"encoding/json"
	"github.com/sujit-baniya/flow"
)

func main() {
	messages := []string{"this", "is", "a", "test"}
	bt, _ := json.Marshal(messages)
	flow.AddNode("for-each-message", ForEachMessage)
	flow.AddNode("get-messages", GetMessages)
	flow.AddNode("message", Message)
	flow.AddNode("send", Send)
	flow1 := flow.New()
	flow1.Loop("for-each-message", "message")
	flow1.Edge("get-messages", "for-each-message")
	flow1.Edge("message", "send")
	_, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: bt,
	})
	if e != nil {
		panic(e)
	}
}

func GetMessages(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

func Message(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

func Send(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

func ForEachMessage(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}
