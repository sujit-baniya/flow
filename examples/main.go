package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID    string
	Name  string
	Phone string
}

func main() {
	users := []User{
		{
			ID:    "1",
			Name:  "Sujit",
			Phone: "+9779856034616",
		},
		{
			ID:    "2",
			Name:  "Baniya",
			Phone: "+9779801634616",
		},
	}
	testThis(users)
}

func testThis(payload Payload) {
	sliceContent := reflect.ValueOf(payload)
	for i := 0; i < sliceContent.Len(); i++ {
		fmt.Println(sliceContent.Index(i).Interface())
	}
}

type Payload interface{}
