package main

import (
	"fmt"
	"reflect"
)


func GenerateBaseCollection(val any) (result []PocketBaseCollection[PbBaseCollectionOptions], err error) {
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

	allCollections := []PocketBaseCollection[PbBaseCollectionOptions]{}
	if err := parseFields(t, &coll, &allCollections); err != nil {
		return result, err
	}
	allCollections = append(allCollections, coll)

	return allCollections, nil
}

func parseFields(t reflect.Type, parentColl *PocketBaseCollection[PbBaseCollectionOptions], allCollections *[]PocketBaseCollection[PbBaseCollectionOptions]) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		a, err := ParseField(field)
		if err != nil {
			if err == ErrStruct {
				childColl := PocketBaseCollection[PbBaseCollectionOptions]{
					Schema: make([]any, 0),
					Name:   field.Type.Name(),
					System: false,
					Type:   "base",
					ID:     GenerateUniqueHash(),
				}
				if err := parseFields(field.Type, &childColl, allCollections); err != nil {
					return err
				}
				*allCollections = append(*allCollections, childColl)

				relationField := PbField[PbRelationFieldOptions]{
					System:      false,
					ID:          GenerateUniqueHash(),
					Name:        field.Name,
					Type:        "relation",
					Required:    false,
					Presentable: false,
					Unique:      false,
					Options: PbRelationFieldOptions{
						CollectionID:  childColl.ID,
						CascadeDelete: false,
						MinSelect:     nil,
						MaxSelect:     1,
						DisplayFields: nil,
					},
				}
				parentColl.Schema = append(parentColl.Schema, relationField)
				continue
			}
			return err
		}
		parentColl.Schema = append(parentColl.Schema, a)
	}
	return nil
}
