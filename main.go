package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Person struct {
	Name  string    `pb:"required,min:1,max:5"`
	Age   int       `pb:"required,min=1,max=10"`
	Date  time.Time `pb:"required,min=2022-01-01 00:00:00.000,max=2022-12-31 23:59:59.000"`
	Inner InnerStruct
}

type InnerStruct struct {
	InnerName string
	InnerAge  int
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
