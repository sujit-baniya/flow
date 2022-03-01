package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/flow"
)

func GetRegistration(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

// VerifyUser Conditional Vertex
func VerifyUser(ctx context.Context, d flow.Data) (flow.Data, error) {
	reg := d.Payload.(Registration)

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

type Registration struct {
	Email    string
	Password string
}

var registeredEmail = map[string]bool{"test@gmail.com": true}

func basicRegistrationFlow() {
	flow1 := flow.New()
	flow1.AddNode("get-registration", GetRegistration)
	flow1.AddNode("create-user", CreateUser)
	flow1.AddNode("cancel-registration", CancelRegistration)
	flow1.AddNode("verify-user", VerifyUser)
	flow1.ConditionalNode("verify-user", map[string]string{
		"pass": "create-user",
		"fail": "cancel-registration",
	})
	flow1.Edge("get-registration", "verify-user")

	registration1 := Registration{
		Email:    "test@gmail.com",
		Password: "admin",
	}
	registration2 := Registration{
		Email:    "test1@gmail.com",
		Password: "admin",
	}
	response, e := flow1.Process(context.Background(), flow.Data{
		Payload: registration1,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
	response, e = flow1.Process(context.Background(), flow.Data{
		Payload: registration2,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func basicRegistrationRawFlow() {
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
	flow1.AddNode("get-registration", GetRegistration)
	flow1.AddNode("create-user", CreateUser)
	flow1.AddNode("cancel-registration", CancelRegistration)
	flow1.AddNode("verify-user", VerifyUser)
	registration1 := Registration{
		Email:    "test@gmail.com",
		Password: "admin",
	}
	registration2 := Registration{
		Email:    "test1@gmail.com",
		Password: "admin",
	}
	response, e := flow1.Process(context.Background(), flow.Data{
		Payload: registration1,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
	response, e = flow1.Process(context.Background(), flow.Data{
		Payload: registration2,
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
