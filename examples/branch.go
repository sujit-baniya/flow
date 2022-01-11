package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sujit-baniya/flow"
)

func GetRegistration(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

// Conditional Vertex
func VerifyUser(ctx context.Context, d flow.Data) (flow.Data, error) {
	var reg Registration
	d.ConvertTo(&reg)
	if _, ok := registeredEmail[reg.Email]; !ok {
		d.Status = "pass"
	} else {
		d.Status = "fail"
	}
	return d, nil
}

func CreateUser(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload(fmt.Sprintf("create user %s", d.Payload))
	return d, nil
}

func CancelRegistration(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload(fmt.Sprintf("cancel user %s", d.Payload))
	return d, nil
}

func registrationNodes() {
	flow.AddNode("get-registration", GetRegistration)
	flow.AddNode("create-user", CreateUser)
	flow.AddNode("cancel-registration", CancelRegistration)
	flow.AddNode("verify-user", VerifyUser)
}

type Registration struct {
	Email    string
	Password string
}

var registeredEmail = map[string]bool{"test@gmail.com": true}

func basicRegistrationFlow() {
	registrationNodes()
	flow1 := flow.New()
	flow1.ConditionalNode("verify-user", map[string]string{
		"pass": "create-user",
		"fail": "cancel-registration",
	})
	flow1.Edge("get-registration", "verify-user")

	registration1 := Registration{
		Email:    "test@gmail.com",
		Password: "admin",
	}
	reg1, _ := json.Marshal(registration1)

	registration2 := Registration{
		Email:    "test1@gmail.com",
		Password: "admin",
	}
	reg2, _ := json.Marshal(registration2)
	response, e := flow1.Process(context.Background(), flow.Data{
		Payload: reg1,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
	response, e = flow1.Process(context.Background(), flow.Data{
		Payload: reg2,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func basicRegistrationRawFlow() {
	registrationNodes()
	rawFlow := []byte(`{
		"edges": [
			["get-registration", "verify-user"]
		],
		"branches":[
			{
				"key": "verify-user",
				"conditional_nodes": {
					"pass": "create-user",
					"fail": "cancel-registration"
				}
			}
		]
	}`)
	flow1 := flow.New(rawFlow)
	registration1 := Registration{
		Email:    "test@gmail.com",
		Password: "admin",
	}
	reg1, _ := json.Marshal(registration1)

	registration2 := Registration{
		Email:    "test1@gmail.com",
		Password: "admin",
	}
	reg2, _ := json.Marshal(registration2)
	response, e := flow1.Process(context.Background(), flow.Data{
		Payload: reg1,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
	response, e = flow1.Process(context.Background(), flow.Data{
		Payload: reg2,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func main() {
	basicRegistrationFlow()
	basicRegistrationRawFlow()
}
