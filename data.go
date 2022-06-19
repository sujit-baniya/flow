package flow

import (
	"encoding/json"
)

type Payload []byte

type Attachment struct {
	Data     []byte
	File     string
	MimeType string
}

type Data struct {
	RequestID    string       `json:"request_id"`
	Payload      Payload      `json:"payload"`
	Status       string       `json:"status"`
	Flow         string       `json:"flow"`
	Operation    string       `json:"operation"`
	FailedReason error        `json:"failed_reason"`
	UserID       uint         `json:"user_id"`
	TimeStamp    int64        `json:"time_stamp"`
	Download     bool         `json:"download"`
	FileName     string       `json:"file_name"`
	Attachments  []Attachment `json:"attachments"`
}

func (d *Data) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Data) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Data) ConvertTo(rs interface{}) error {
	return json.Unmarshal(d.Payload, rs)
}

func (d *Data) ToString() string {
	return string(d.Payload)
}

func (d *Data) GetStatus() string {
	return d.Status
}

func (d *Data) Log() error {
	return nil
}

func (d *Data) LogRecords(count ...int64) error {

	return nil
}

func (d *Data) logUserCount(prefix, date string, count ...int64) error {
	return nil
}

func (d *Data) logUserFlowCount(prefix, date string, count ...int64) error {
	return nil
}

func (d *Data) logUserFlowOperationCount(prefix, date string) error {
	return nil
}

func (d *Data) logUserFlowOperation(prefix string) error {
	return nil
}
