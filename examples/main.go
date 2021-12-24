package main

import (
	"fmt"
	"github.com/sujit-baniya/flow"
)

func main() {
	flow1 := flow.New()
	flow1.Key = "process-avatar-upload"
	flow1.Node("get-image", GetImage)
	flow1.Node("create-avatar", CreateAvatar)
	flow1.Node("delete-image", DeleteImage)
	flow1.Node("send-email", SendEmail)
	flow1.ConditionalNode("handle-image", HandleImage, map[string]string{
		"pass": "create-avatar",
		"fail": "delete-image",
	})
	flow1.Edge("get-image", "handle-image")
	flow1.Edge("handle-image", "send-email")
	r, _ := flow1.Process(flow.DataSource{Payload: flow.Payload("https://sujitbaniya.com/logo.svg")})
	fmt.Println(string(r.Payload))
}

func data() flow.Payload {
	return flow.Payload(`{
	"nodes": ["send-single", "get-provider", "send-message", "store-message", "send-callback"],
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
