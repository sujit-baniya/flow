package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sujit-baniya/flow"
	"strings"
)

func GetSentence(ctx context.Context, d flow.Data) (flow.Data, error) {
	words := strings.Split(d.ToString(), ` `)
	bt, _ := json.Marshal(words)
	d.Payload = bt
	return d, nil
}

func ForEachWord(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

func WordUpperCase(ctx context.Context, d flow.Data) (flow.Data, error) {
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	d.Payload = flow.Payload(strings.ToTitle(strings.ToLower(word)))
	return d, nil
}

func AppendString(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload("Upper Case: " + string(d.Payload))
	fmt.Println(d.ToString())
	return d, nil
}

func wordNodes(flow1 *flow.Flow) {
	flow1.AddNode("get-sentence", GetSentence)
	flow1.AddNode("for-each-word", ForEachWord)
	flow1.AddNode("upper-case", WordUpperCase)
	flow1.AddNode("append-string", AppendString)
}

func main() {
	flow1 := flow.New()
	wordNodes(flow1)
	flow1.Loop("for-each-word", "upper-case")
	flow1.Edge("get-sentence", "for-each-word")
	flow1.Edge("upper-case", "append-string")
	_, e := flow1.Process(context.Background(), flow.Data{
		Payload: flow.Payload("this is a sentence"),
	})
	if e != nil {
		panic(e)
	}
}
