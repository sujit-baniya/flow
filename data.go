package flow

import (
	"github.com/goccy/go-json"
)

type Payload []byte

type Data struct {
	RequestID     string  `json:"request_id"`
	Payload       Payload `json:"payload"`
	Status        string  `json:"status"`
	CurrentVertex string  `json:"current_vertex"`
	FailedReason  error   `json:"failed_reason"`
}

func (d Data) ConvertTo(rs interface{}) error {
	return json.Unmarshal(d.Payload, rs)
}

func (d Data) ToString() string {
	return string(d.Payload)
}

func (d Data) GetStatus() string {
	return d.Status
}
