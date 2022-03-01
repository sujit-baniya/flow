package flow

type Payload interface{}

type Data struct {
	RequestID     string  `json:"request_id"`
	Payload       Payload `json:"payload"`
	Status        string  `json:"status"`
	CurrentVertex string  `json:"current_vertex"`
	FailedReason  error   `json:"failed_reason"`
}

func (d Data) GetStatus() string {
	return d.Status
}
