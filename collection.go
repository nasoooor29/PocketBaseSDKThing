package main

import (
	"fmt"
	"reflect"
)

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
	coll.System = false
	coll.ID = GenerateUniqueHash()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		a, err := ParseField(field)
		if err != nil {
			return result, err
		}
		coll.Schema = append(coll.Schema, a)
	}
	return coll, nil
}
