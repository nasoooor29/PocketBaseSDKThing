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

func ParseTextField(f reflect.StructField) (opts PbField[PbTextFieldOptions], err error) {
	fields := GetTheFields(f)
	minStr, found := fields["min"]
	if found {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return opts, fmt.Errorf("%v min must be a number on text field", f.Name)
		}
		opts.Options.Min = min
	}
	maxStr, found := fields["max"]
	if found {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return opts, fmt.Errorf("%v max must be a number on text field", f.Name)
		}
		opts.Options.Max = max
	}
	pattern, found := fields["pattern"]
	if found {
		_, err = regexp.Compile(pattern)
		if err != nil {
			return opts, fmt.Errorf("%v pattern must be a valid regex on text field", f.Name)
		}
		opts.Options.Pattern = pattern
	}
	opts.ID = GenerateUniqueHash()
	opts.Name = f.Name
	_, opts.Required = fields["required"]
	_, opts.Presentable = fields["presentable"]
	_, opts.Unique = fields["unique"]
	opts.Type = "text"
	return opts, nil

}

func ParseNumberOptions(f reflect.StructField) (opts PbField[PbNumberFieldOptions], err error) {
	fields := GetTheFields(f)
	minStr, found := fields["min"]
	if found {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return opts, fmt.Errorf("%v min must be a number on number field", f.Name)
		}
		opts.Options.Min = min
	}
	maxStr, found := fields["max"]
	if found {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return opts, fmt.Errorf("%v max must be a number on number field", f.Name)
		}
		opts.Options.Max = max
	}
	opts.ID = GenerateUniqueHash()
	opts.Name = f.Name
	_, opts.Required = fields["required"]
	_, opts.Presentable = fields["presentable"]
	_, opts.Unique = fields["unique"]
	_, opts.Options.NoDecimal = fields["nodecimal"]
	opts.Type = "number"
	return opts, nil
}

func ParseBoolOptions(f reflect.StructField) (opts PbField[PbBoolFieldOptions], err error) {
	opts.ID = GenerateUniqueHash()
	fields := GetTheFields(f)
	opts.Name = f.Name
	_, opts.Required = fields["required"]
	_, opts.Presentable = fields["presentable"]
	_, opts.Unique = fields["unique"]
	opts.Type = "boolean"
	return opts, nil
}

func ParseField(f reflect.StructField) (any, error) {
	t := f.Tag.Get("type")
	if t != "" {
		return nil, fmt.Errorf("type %v not supported", f.Type.Kind())
	}
	switch f.Type.Kind() {
	case reflect.String:
		return ParseTextField(f)
	case reflect.Int:
		return ParseNumberOptions(f)
	case reflect.Bool:
		return ParseBoolOptions(f)
	case reflect.Ptr: // this need to be recursive
		return nil, fmt.Errorf("pointer type not supported")
	case reflect.Struct:
		return nil, fmt.Errorf("struct type not supported")
	default:
		return nil, fmt.Errorf("type %v not supported", f.Type.Kind())
	}
}
