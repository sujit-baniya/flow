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

func basicNodes(flow1 *flow.Flow) {
	flow1.AddNode("message", Message)
	flow1.AddNode("send", Send)
}

func basicFlow() {
	flow1 := flow.New()
	basicNodes(flow1)
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
	rawFlow := []byte(`{
		"edges": [
			["message", "send"]
		]
	}`)
	flow1 := flow.New(rawFlow)
	basicNodes(flow1)
	response, e := flow1.Process(context.Background(), flow.Data{
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
