package flow

import "context"

type Node interface {
	Process(ctx context.Context, data Data) (Data, error)
	AddEdge(node Node)
	IsBranch() bool
	GetKey() string
}

type Payload []byte

type Handler func(ctx context.Context, data Data) (Data, error)

var NodeList = map[string]Handler{}

func AddNode(node string, handler Handler) {
	NodeList[node] = handler
}

func GetNodeHandler(node string) Handler {
	return NodeList[node]
}
