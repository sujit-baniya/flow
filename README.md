# Introduction
This package provides simple graph to execute functions in a group.

## Features
- Define nodes of different types: Vertex, Branch and Loop
- Define branch for conditional nodes


## Installation
> go get github.com/sujit-baniya/flow

## Usage

### Basic Flow
```go
package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/flow"
)

func Message(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload(fmt.Sprintf("message %s", d.Payload))
	return d, nil
}

func Send(ctx context.Context, d flow.Data) (flow.Data, error) {
	d.Payload = flow.Payload(fmt.Sprintf("This is send %s", d.Payload))
	return d, nil
}

func nodes() {
	flow.AddNode("message", Message)
	flow.AddNode("send", Send)
}

func basicFlow() {
	nodes()
	flow1 := flow.New()
	flow1.Edge("message", "send")
	response, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: flow.Payload("Payload"),
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func basicRawFlow() {
	nodes()
	rawFlow := []byte(`{
		"edges": [
			["message", "send"]
		]
	}`)
	flow1 := flow.New(rawFlow)
	response, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: flow.Payload("Payload"),
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
}

func main() {
	basicFlow()
	basicRawFlow()
}

```

### Branch Flow
```go
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
	flow.AddBranch("verify-user", VerifyUser)
	flow.AddNode("create-user", CreateUser)
	flow.AddNode("cancel-registration", CancelRegistration)
}

type Registration struct {
	Email string
	Password string
}

var registeredEmail = map[string]bool {"test@gmail.com": true}

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
	response, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: reg1,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
	response, e = flow1.Build().Process(context.Background(), flow.Data{
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
	response, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: reg1,
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(response.ToString())
	response, e = flow1.Build().Process(context.Background(), flow.Data{
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

```

### Loop Flow
```go
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
	d.Payload = flow.Payload(strings.ToTitle(d.ToString()))
	fmt.Println(d.ToString())
	return d, nil
}

func wordNodes() {
	flow.AddNode("get-sentence", GetSentence)
	flow.AddNode("for-each-word", ForEachWord)
	flow.AddNode("upper-case", WordUpperCase)
}

func main() {
	wordNodes()
	flow1 := flow.New()
	flow1.Loop("for-each-word", "upper-case")
	flow1.Edge("get-sentence", "for-each-word")
	_, e := flow1.Build().Process(context.Background(), flow.Data{
		Payload: flow.Payload("this is a sentence"),
	})
	if e != nil {
		panic(e)
	}
}

```
## ToDo List
- Implement async flow and nodes
- Implement distributed nodes
