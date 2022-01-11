package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/flow"
)

func Message(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload(fmt.Sprintf("message %s", d.Payload))
	return d, nil
}

func Send(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload(fmt.Sprintf("This is send %s", d.Payload))
	return d, nil
}

func basicNodes() {
	flow.AddNode("message", Message)
	flow.AddNode("send", Send)
}

func basicFlow() {
	basicNodes()
	flow1 := flow.New()
	flow1.Edge("message", "send")
	response, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: flow.Payload("Payload"),
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func basicRawFlow() {
	basicNodes()
	rawFlow := []byte(`{
		"edges": [
			["message", "send"]
		]
	}`)
	flow1 := flow.New(rawFlow)
	response, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: flow.Payload("Payload"),
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func main() {
	basicFlow()
	basicRawFlow()
}
