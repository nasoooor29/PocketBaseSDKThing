package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetTagsMap(f reflect.StructField) map[string]string {
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
		// field = strings.Split(val, " ")
		// if len(field) == 2 {
		// 	fields[strings.ToLower(field[0])] = field[1]
		// }
		fields[strings.ToLower(val)] = ""
	}
	return fields
}

func ParseTextField(f reflect.StructField) (opts PbField[PbTextFieldOptions], err error) {
	fields := GetTagsMap(f)
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
	fields := GetTagsMap(f)
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
	fields := GetTagsMap(f)
	opts.Name = f.Name
	_, opts.Required = fields["required"]
	_, opts.Presentable = fields["presentable"]
	_, opts.Unique = fields["unique"]
	opts.Type = "boolean"
	return opts, nil
}

// TODO: parse min and max for date (2024-07-20 12:00:00.000Z format)
func ParseDateField(f reflect.StructField) (opts PbField[PbDateFieldOptions], err error) {
	fields := GetTagsMap(f)
	min, found := fields["min"]
	if found {
		fmt.Printf("min: %v\n", min)
		time, err := time.Parse("2006-01-02 15:04:05.000", min)
		if err != nil {
			return opts, fmt.Errorf("%v min must be a valid date on date field eg:2006-01-02 15:04:05.000", f.Name)
		}
		opts.Options.Min = time.String()
	}
	max, found := fields["max"]
	if found {
		time, err := time.Parse("2006-01-02 15:04:05.000", max)
		if err != nil {
			return opts, fmt.Errorf("%v max must be a valid date on date field eg:2006-01-02 15:04:05.000", f.Name)
		}
		opts.Options.Max = time.String()
	}
	opts.ID = GenerateUniqueHash()
	opts.Name = f.Name
	_, opts.Required = fields["required"]
	_, opts.Presentable = fields["presentable"]
	_, opts.Unique = fields["unique"]
	opts.Type = "date"
	return opts, nil
}

func ParseField(f reflect.StructField) (any, error) {
	// ? available types: text,file,relation,editor,number,bool,email,url,date,select,json
	// ? need to be done types: file,relation,email,url,select,json

	fields := GetTagsMap(f)
	v, found := fields["type"]
	if found && v == "editor" {
		v, err := ParseTextField(f)
		if err != nil {
			return nil, err
		}
		v.Type = "editor"
		return v, nil
	}
	if (f.Type.Kind() == reflect.Ptr || f.Type.Kind() == reflect.Struct) && f.Type.PkgPath() == "time" {
		return ParseDateField(f)
	}
	switch f.Type.Kind() {
	case reflect.String:
		return ParseTextField(f)
	case reflect.Int:
		return ParseNumberOptions(f)
	case reflect.Bool:
		return ParseBoolOptions(f)
	case reflect.Struct:
		return nil, ErrStruct
	case reflect.Ptr:
		return nil, ErrStruct
	default:
		return nil, fmt.Errorf("type %v not supported on %v", f.Type.Kind(), f.Name)
	}
}
