package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MyCustomStruct struct {
	Name        string       `pb:"min=1,max=100"` // for these i think we should use tags to define the options
	Age         int          `pb:"min=1,max=300,NoDecimal"`
	Active      bool         `pb:"required"`
	CustomField string       `json:"customField"`
	nest        NestedStruct // if pointer i think we should make another table and make a relation
	// if no pointer we add the fields on the parent struct
}

type NestedStruct struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Active bool   `json:"active"`
}

func main() {
	myStruct := MyCustomStruct{
		Name:        "test",
		Age:         32,
		Active:      true,
		CustomField: "a1",
		nest: NestedStruct{
			Name:   "a1",
			Age:    32,
			Active: false,
		},
	}

	j, err := GenerateBaseCollection(myStruct)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	data, err := json.Marshal(j)
	if err != nil {
		return
	}
	err = os.WriteFile("test.json", data, 0644)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}
