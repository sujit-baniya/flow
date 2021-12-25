package main

import (
	"fmt"
	"github.com/sujit-baniya/flow"
)

func SendSingle(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func EstimateSingle(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func DeductBalance(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func GetProvider(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func SendMessage(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func StoreMessage(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func SendCallback(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func nodes() {
	flow.AddNode("send-single", SendSingle)
	flow.AddNode("estimate-single", EstimateSingle)
	flow.AddNode("deduct-balance", DeductBalance)
	flow.AddNode("get-provider", GetProvider)
	flow.AddNode("send-message", SendMessage)
	flow.AddNode("store-message", StoreMessage)
	flow.AddNode("send-callback", SendCallback)
}

func branches() {
	flow.AddBranch("check-sender-id", CheckSenderID)
	flow.AddBranch("check-message", CheckMessage)
	flow.AddBranch("prepare-message", PrepareMessage)
	flow.AddBranch("check-balance", CheckBalance)
}

func CheckSenderID(source flow.DataSource) (flow.DataSource, error) {
	source.Status = "pass"
	return source, nil
}

func CheckMessage(source flow.DataSource) (flow.DataSource, error) {
	source.Status = "pass"
	return source, nil
}

func PrepareMessage(source flow.DataSource) (flow.DataSource, error) {
	source.Status = "success"
	return source, nil
}

func CheckBalance(source flow.DataSource) (flow.DataSource, error) {
	source.Status = "pass"
	return source, nil
}

func main() {
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
	flow1.Edge("send-single", "check-sender-id")
	flow1.Edge("check-sender-id", "check-message")
	flow1.Edge("estimate-single", "check-balance")
	flow1.Edge("deduct-balance", "prepare-message")
	flow1.Edge("get-provider", "send-message")
	flow1.Edge("send-message", "store-message")
	flow1.Edge("store-message", "send-callback")
	res, err := flow1.Build().Process(flow.DataSource{
		Payload:   flow.Payload(`{"email": "s.baniy8a.np@gmail.com", "password": "123456", "avatar": "image.svg"}`),
		RequestID: "asdasdas",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToString())
}

func data() flow.Payload {
	return flow.Payload(`{
	"nodes": ["send-single", "estimate-single", "deduct-balance", "get-provider", "send-message", "store-message", "send-callback"],
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
		["check-message", "estimate-single"],
		["estimate-single", "check-balance"],
		["check-balance", "deduct-balance"],
		["deduct-balance", "prepare-message"],
		["prepare-message", "get-provider"],
		["get-provider", "send-message"],
		["send-message", "store-message"],
		["store-message", "send-callback"]
	]
}`)
}
