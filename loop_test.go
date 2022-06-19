package flow

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func GetSentence(ctx context.Context, d Data) (Data, error) {
	words := strings.Split(d.ToString(), ` `)
	bt, _ := json.Marshal(words)
	d.Payload = bt
	return d, nil
}

func ForEachWord(ctx context.Context, d Data) (Data, error) {
	return d, nil
}

func WordUpperCase(ctx context.Context, d Data) (Data, error) {
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	bt, _ := json.Marshal(strings.Title(strings.ToLower(word)))
	d.Payload = bt
	return d, nil
}

func AppendIP(ctx context.Context, d Data) (Data, error) {
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	bt, _ := json.Marshal("IP: " + word)
	d.Payload = bt
	return d, nil
}

func AppendString(ctx context.Context, d Data) (Data, error) {
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	bt, _ := json.Marshal("Upper Case: " + word)
	d.Payload = bt
	return d, nil
}

func BenchmarkFlow_Loop(b *testing.B) {
	flow1 := New()
	flow1.AddNode("get-sentence", GetSentence)
	flow1.AddNode("for-each-word", ForEachWord)
	flow1.AddNode("upper-case", WordUpperCase)
	flow1.AddNode("append-string", AppendString)
	flow1.AddNode("append-ip", AppendIP)
	flow1.Loop("for-each-word", "append-ip", "upper-case")
	flow1.Edge("get-sentence", "for-each-word")
	flow1.Edge("upper-case", "append-string")
	for i := 0; i < b.N; i++ {
		flow1.Process(context.Background(), Data{
			Payload: Payload("this is a sentence"),
		})
	}
}
