package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func GetTheFields(f reflect.StructField) map[string]string {
	pbType, found := f.Tag.Lookup("pb")
	if !found {
		return nil
	}
	values := strings.Split(pbType, ",")
	fields := make(map[string]string)
	for _, val := range values {
		field := strings.Split(val, ":")
		if len(field) == 2 {
			fields[strings.ToLower(field[0])] = field[1]
		}
		field = strings.Split(val, "=")
		if len(field) == 2 {
			fields[strings.ToLower(field[0])] = field[1]
		}
		field = strings.Split(val, " ")
		if len(field) == 2 {
			fields[strings.ToLower(field[0])] = field[1]
		}
		fields[strings.ToLower(val)] = ""
	}
	return fields
}

func GenerateFieldData(f reflect.StructField) (data FieldData) {
	data.ID = GenerateUniqueHash()
	data.Name = f.Name
	fields := GetTheFields(f)

	_, found := fields["required"]
	if found {
		data.Required = true
	}
	_, found = fields["presentable"]
	if found {
		data.Presentable = true
	}
	_, found = fields["unique"]
	if found {
		data.Unique = true
	}

	return data
}

func ParseTextOptions(f reflect.StructField) (opts PbTextFieldOptions, err error) {
	fields := GetTheFields(f)
	minStr, found := fields["min"]
	if found {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return opts, fmt.Errorf("%v min must be a number on text field", f.Name)
		}
		opts.Min = min
	}
	maxStr, found := fields["max"]
	if found {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return opts, fmt.Errorf("%v max must be a number on text field", f.Name)
		}
		opts.Max = max
	}
	pattern, found := fields["pattern"]
	if found {
		_, err = regexp.Compile(pattern)
		if err != nil {
			return opts, fmt.Errorf("%v pattern must be a valid regex on text field", f.Name)
		}
		opts.Pattern = pattern
	}
	return opts, nil

}

func ParseNumberOptions(f reflect.StructField) (opts PbNumberFieldOptions, err error) {
	fields := GetTheFields(f)
	minStr, found := fields["min"]
	if found {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return opts, fmt.Errorf("%v min must be a number on number field", f.Name)
		}
		opts.Min = min
	}
	maxStr, found := fields["max"]
	if found {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return opts, fmt.Errorf("%v max must be a number on number field", f.Name)
		}
		opts.Max = max
	}
	_, found = fields["nodecimal"]
	if found {
		opts.NoDecimal = true
	}
	return opts, nil

}

func ParseField(f reflect.StructField) (any, error) {
	t := f.Tag.Get("type")
	data := GenerateFieldData(f)
	if t == "" && f.Type.Kind() == reflect.String {
		opts, err := ParseTextOptions(f)
		if err != nil {
			return nil, err
		}
		result := &PbField[PbTextFieldOptions]{
			FieldData: &data,
			Options:   opts,
		}
		result.FieldData.Type = "text"
		return result, nil
	}
	if t == "" && f.Type.Kind() == reflect.Int {
		opts, err := ParseNumberOptions(f)
		if err != nil {
			return nil, err
		}
		result := &PbField[PbNumberFieldOptions]{
			FieldData: &data,
			Options:   opts,
		}
		result.FieldData.Type = "number"
		return result, nil
	}
	if t == "" && f.Type.Kind() == reflect.Bool {
		result := &PbField[PbBoolFieldOptions]{
			FieldData: &data,
		}
		result.FieldData.Type = "boolean"
		return result, nil
	}
	if t == "" && f.Type.Kind() == reflect.Ptr { // this need to be recursive
		result := &PbField[PbRelationFieldOptions]{
			FieldData: &data,
		}
		result.FieldData.Type = "relation"
		return result, nil // deal with the relation later
	}
	if t == "" && f.Type.Kind() == reflect.Struct {
		// result := &PbField[PbObjectFieldOptions]{
		// 	FieldData: &data,
		// }
		// result.FieldData.Type = "object"
		// return result, nil // deal with the relation
		return nil, fmt.Errorf("struct type not supported")
	}
	return nil, fmt.Errorf("type %v not supported", f.Type.Kind())
}
