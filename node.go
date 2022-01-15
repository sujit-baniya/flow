package flow

import "context"

type Node interface {
	Process(ctx context.Context, data Data) (Data, error)
	AddEdge(node Node)
	GetType() string
	GetKey() string
}

type Handler func(ctx context.Context, data Data) (Data, error)

var nodeList = map[string]Handler{}

func AddNode(node string, handler Handler) {
	nodeList[node] = handler
}

func GetNodeHandler(node string) Handler {
	return nodeList[node]
}
