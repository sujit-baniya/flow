package flow

import (
	"context"
	"strings"
	"testing"
)

func GetSentence(ctx context.Context, d Data) (Data, error) {
	words := strings.Split(d.Payload.(string), ` `)
	d.Payload = words
	return d, nil
}

func ForEachWord(ctx context.Context, d Data) (Data, error) {
	return d, nil
}

func WordUpperCase(ctx context.Context, d Data) (Data, error) {
	d.Payload = strings.Title(strings.ToLower(d.Payload.(string)))
	return d, nil
}

func AppendIP(ctx context.Context, d Data) (Data, error) {
	var s strings.Builder
	s.WriteString("IP: ")
	s.WriteString(d.Payload.(string))
	d.Payload = s.String()
	return d, nil
}

func AppendString(ctx context.Context, d Data) (Data, error) {
	var s strings.Builder
	s.WriteString("Append: ")
	s.WriteString(d.Payload.(string))
	d.Payload = s.String()
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
