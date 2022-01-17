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

func basicFlow() {
	flow1 := flow.New()
	flow1.AddNode("message", Message)
	flow1.AddNode("send", Send)
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
	rawFlow := []byte(`{
		"edges": [
			["message", "send"]
		]
	}`)
	flow1 := flow.New(rawFlow)
	flow1.AddNode("message", Message)
	flow1.AddNode("send", Send)
	response, e := flow1.Process(context.Background(), flow.Data{
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
	"fmt"
	"encoding/json"
	"github.com/sujit-baniya/flow"
)

func GetRegistration(ctx context.Context, d flow.Data) (flow.Data, error) {
	return d, nil
}

// VerifyUser Conditional Vertex
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

```

### Loop Flow
```go
package main

import (
	"context"
	"fmt"
	"encoding/json"
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
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	bt, _ := json.Marshal(strings.Title(strings.ToLower(word)))
	d.Payload = bt
	return d, nil
}

func AppendString(ctx context.Context, d flow.Data) (flow.Data, error) {
	var word string
	_ = json.Unmarshal(d.Payload, &word)
	bt, _ := json.Marshal("Upper Case: " + word)
	d.Payload = bt
	return d, nil
}

func main() {
	flow1 := flow.New()
	flow1.AddNode("get-sentence", GetSentence)
	flow1.AddNode("for-each-word", ForEachWord)
	flow1.AddNode("upper-case", WordUpperCase)
	flow1.AddNode("append-string", AppendString)
	flow1.Loop("for-each-word", "upper-case")
	flow1.Edge("get-sentence", "for-each-word")
	flow1.Edge("upper-case", "append-string")
	resp, e := flow1.Process(context.Background(), flow.Data{
		Payload: flow.Payload("this is a sentence"),
	})
	if e != nil {
		panic(e)
	}
	fmt.Println(resp.ToString())
}

```
## ToDo List
- Implement async flow and nodes
- Implement distributed nodes
