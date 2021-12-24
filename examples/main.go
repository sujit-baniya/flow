package main

import (
	"github.com/sujit-baniya/flow"
)

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
