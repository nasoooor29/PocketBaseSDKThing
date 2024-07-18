package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string `pb:"required,min:1,max:5"`
	Age  int `pb:"required,min=1,max=10"`
}

func main() {

	j, err := GenerateBaseCollection(
		Person{
			Name: "",
			Age:  0,
		},
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	data, err := json.Marshal([]any{j})
	if err != nil {
		return
	}
	err = os.WriteFile("test.json", data, 0644)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}
