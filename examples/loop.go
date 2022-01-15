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
	bt, _ := json.Marshal(strings.Title(strings.ToLower(word)))
	d.Payload = bt
	return d, nil
}

func AppendString(ctx context.Context, d flow.Data) (flow.Data, error) {
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	bt, _ := json.Marshal("Upper Case: " + word)
	d.Payload = bt
	return d, nil
}

func wordNodes() {
	flow.AddNode("get-sentence", GetSentence)
	flow.AddNode("for-each-word", ForEachWord)
	flow.AddNode("upper-case", WordUpperCase)
	flow.AddNode("append-string", AppendString)
}

func main() {
	wordNodes()
	flow1 := flow.New()
	flow1.Loop("for-each-word", "upper-case")
	flow1.Edge("get-sentence", "for-each-word")
	flow1.Edge("upper-case", "append-string")
	resp, e := flow1.Process(context.Background(), flow.Data{
		Payload: flow.Payload("this is a sentence"),
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(resp.ToString())
}
