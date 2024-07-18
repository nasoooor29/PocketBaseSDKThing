package main

import (
	"fmt"
	"reflect"
)

// input:
type MyCustomStruct struct {
	Name  string
	Age   int
	Inner MyInnerStruct
}

type MyInnerStruct struct {
	InnerName string
	InnerAge  int
}

// output:
type ResultStruct struct {
	Name      string
	Age       int
	InnerName string
	InnerAge  int
}

// i know it should be like but for "simplicity" func GenerateCollectionSchema[T CollectionType](val any) (result PocketBaseCollection[T], err error) {
func GenerateBaseCollection(val any) (result PocketBaseCollection[PbBaseCollectionOptions], err error) {
	t := reflect.TypeOf(val)
	if t.Kind() != reflect.Struct && t.Kind() != reflect.Map {
		return result, fmt.Errorf("val must be a struct or a map")
	}
	coll := PocketBaseCollection[PbBaseCollectionOptions]{}
	coll.Schema = make([]any, 0)
	coll.Name = t.Name()
	coll.System = false
	coll.Type = "base"
	coll.ID = GenerateUniqueHash()
	if err := parseFields(t, &coll); err != nil {
		return result, err
	}
	return coll, nil
}

func parseFields(t reflect.Type, coll *PocketBaseCollection[PbBaseCollectionOptions]) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("t.Name: %v, %v.%v\n", t.Name(), t.Name(), field.Name)
		a, err := ParseField(field)
		if err == nil {
			coll.Schema = append(coll.Schema, a)
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			if err := parseFields(field.Type, coll); err != nil {
				return err
			}
			continue
		}
		return err
	}
	return nil
}
