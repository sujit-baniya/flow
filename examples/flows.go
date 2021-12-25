package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sujit-baniya/flow"
)

func registerNodes() {
	flow.AddNode("get-data", GetData)
	flow.AddNode("store-data", StoreData)
	flow.AddNode("cancel-registration", CancelRegistration)
	flow.AddNode("get-image", GetImage)
	flow.AddNode("create-avatar", CreateAvatar)
	flow.AddNode("delete-image", DeleteImage)
	flow.AddNode("throw-error", ThrowError)

	flow.AddBranch("handle-registration", HandleRegistration)
	flow.AddBranch("handle-image", HandleImage)
}

func registrationFlow() {
	flow1 := flow.New()
	flow1.Key = "process-registration"

}

func main() {
	registerNodes()
	flow1 := flow.New()
	flow1.ConditionalNode("handle-registration", map[string]string{
		"pass": "get-image",
		"fail": "cancel-registration",
	})
	flow1.ConditionalNode("handle-image", map[string]string{
		"pass": "create-avatar",
		"fail": "delete-image",
	})
	flow1.Edge("get-data", "handle-registration")
	flow1.Edge("get-image", "handle-image")
	flow1.Edge("cancel-registration", "throw-error")
	res, err := flow1.Build().Process(flow.DataSource{
		Payload:   flow.Payload(`{"email": "s.baniy8a.np@gmail.com", "password": "123456", "avatar": "image.svg"}`),
		RequestID: "asdasdas",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.ToString())
}

type Registration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

func ThrowError(data flow.DataSource) (flow.DataSource, error) {
	return data, data.FailedReason
}

func GetData(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func CancelRegistration(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func StoreData(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func GetImage(data flow.DataSource) (flow.DataSource, error) {
	return data, nil
}

func HandleRegistration(data flow.DataSource) (flow.DataSource, error) {
	var registration Registration
	data.ConvertTo(&registration)
	data.Status = "fail"
	if registration.Email == "s.baniya.np@gmail.com" && registration.Password == "123456" {
		data.Status = "pass"
	} else {
		data.FailedReason = errors.New("Invalid Credentials")
	}
	return data, nil
}

func HandleImage(data flow.DataSource) (flow.DataSource, error) {
	var registration Registration
	data.ConvertTo(&registration)

	status := "fail"
	if registration.Avatar != "" {
		status = "pass"
	}
	data.Status = status
	return data, nil
}

func CreateAvatar(data flow.DataSource) (flow.DataSource, error) {
	data.Payload = flow.Payload(fmt.Sprintf("creating avatar from image %s", data))
	return data, nil
}

func DeleteImage(data flow.DataSource) (flow.DataSource, error) {
	data.Payload = flow.Payload(fmt.Sprintf("delete image %s", data))
	return data, nil
}

func SendEmail(data flow.DataSource) (flow.DataSource, error) {
	data.Payload = flow.Payload(fmt.Sprintf("I'm sending email for %s", data))
	return data, nil
}
