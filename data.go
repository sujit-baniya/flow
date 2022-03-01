package flow

import "encoding/json"

type Payload interface{}

type Data struct {
	RequestID     string  `json:"request_id"`
	Payload       Payload `json:"payload"`
	Status        string  `json:"status"`
	CurrentVertex string  `json:"current_vertex"`
	FailedReason  error   `json:"failed_reason"`
}

func (d Data) ConvertTo(rs interface{}) error {
	switch v := d.Payload.(type) {
	case []byte:
		return json.Unmarshal(v, rs)
	}
	return nil
}

func (d Data) ToString() string {
	switch v := d.Payload.(type) {
	case []byte:
		return string(v)
	}
	return d.Payload.(string)
}

func (d Data) GetStatus() string {
	return d.Status
}
