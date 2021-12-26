package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sujit-baniya/flow"
)

func SendSingle(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func EstimateSingle(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func DeductBalance(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func GetProvider(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func SendMessage(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func StoreMessage(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func SendCallback(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func ReceiveRequest(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, nil
}

func ThrowError(ctx context.Context, data flow.Data) (flow.Data, error) {
	return data, data.FailedReason
}

func ForEachMessage(ctx context.Context, d flow.Data) (flow.Data, error) {

	return d, nil
}

func nodes() {
	flow.AddNode("for-each-message", ForEachMessage)
	flow.AddNode("receive-request", ReceiveRequest)
	flow.AddNode("send-single", SendSingle)
	flow.AddNode("estimate-single", EstimateSingle)
	flow.AddNode("deduct-balance", DeductBalance)
	flow.AddNode("get-provider", GetProvider)
	flow.AddNode("send-message", SendMessage)
	flow.AddNode("store-message", StoreMessage)
	flow.AddNode("send-callback", SendCallback)
	flow.AddNode("throw-error", ThrowError)
}

func branches() {
	flow.AddBranch("check-sender-id", CheckSenderID)
	flow.AddBranch("check-message", CheckMessage)
	flow.AddBranch("prepare-message", PrepareMessage)
	flow.AddBranch("check-balance", CheckBalance)
	flow.AddBranch("validate-request", ValidateRequest)
}

func CheckSenderID(ctx context.Context, source flow.Data) (flow.Data, error) {
	source.Status = "pass"
	return source, nil
}

func ValidateRequest(ctx context.Context, source flow.Data) (flow.Data, error) {
	source.Status = "pass"
	return source, nil
}

func CheckMessage(ctx context.Context, source flow.Data) (flow.Data, error) {
	source.Status = "pass"
	return source, nil
}

func PrepareMessage(ctx context.Context, source flow.Data) (flow.Data, error) {
	source.Status = "success"
	return source, nil
}

func CheckBalance(ctx context.Context, source flow.Data) (flow.Data, error) {
	source.Status = "pass"
	return source, nil
}

func main() {
	rawFlow()
	// normalFlow()
}

func rawFlow() {
	nodes()
	branches()
	flow1 := flow.New(data())
	res, err := flow1.Build().Process(context.Background(), flow.Data{
		Payload:   flow.Payload(`{"email": "s.baniy8a.np@gmail.com", "password": "123456", "avatar": "image.svg"}`),
		RequestID: "asdasdas",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToString())
}

func normalFlow() {
	nodes()
	branches()
	flow1 := flow.New()
	flow1.ConditionalNode("prepare-message", map[string]string{
		"error":        "store-message",
		"success":      "get-provider",
		"unsubscribed": "store-message",
		"invalid":      "store-message",
	})
	flow1.ConditionalNode("check-balance", map[string]string{
		"pass": "deduct-balance",
		"fail": "prepare-message",
	})
	flow1.ConditionalNode("check-sender-id", map[string]string{
		"pass": "check-message",
		"fail": "prepare-message",
	})
	flow1.ConditionalNode("check-message", map[string]string{
		"pass": "estimate-single",
		"fail": "prepare-message",
	})
	flow1.ConditionalNode("validate-request", map[string]string{
		"pass": "for-each-message",
		"fail": "throw-error",
	})
	flow1.Edge("receive-request", "validate-request")
	flow1.Loop("for-each-message", "send-single")
	flow1.Edge("send-single", "check-sender-id")
	flow1.Edge("check-sender-id", "check-message")
	flow1.Edge("estimate-single", "check-balance")
	flow1.Edge("deduct-balance", "prepare-message")
	flow1.Edge("get-provider", "send-message")
	flow1.Edge("send-message", "store-message")
	flow1.Edge("store-message", "send-callback")
	messages := []string{"this", "is", "a", "test"}
	bt, _ := json.Marshal(messages)
	res, err := flow1.Build().Process(context.Background(), flow.Data{
		Payload:   bt,
		RequestID: "asdasdas",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToString())
}

func data() flow.Payload {
	return flow.Payload(`{
	"branches":[
		{
			"key": "check-sender-id",
			"conditional_nodes": {
				"pass": "check-message",
				"fail": "prepare-message"
			}
		},
		{
			"key": "check-message",
			"conditional_nodes": {
				"pass": "estimate-single",
				"fail": "prepare-message"
			}
		},
		{
			"key": "check-balance",
			"conditional_nodes": {
				"pass": "deduct-balance",
				"fail": "prepare-message"
			}
		},
		{
			"key": "prepare-message",
			"conditional_nodes": {
				"error": "store-message",
				"success": "get-provider",
				"unsubscribed": "store-message",
				"invalid": "store-message"
			}
		}
	],
	"edges": [
		["send-single", "check-sender-id"],
		["check-sender-id", "check-message"],
		["estimate-single", "check-balance"],
		["deduct-balance", "prepare-message"],
		["get-provider", "send-message"],
		["send-message", "store-message"],
		["store-message", "send-callback"]
	]
}`)
}
