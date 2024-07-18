package main

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"time"
)

func GenerateUniqueHash() string {
	timeBytes := []byte(time.Now().Format("2006-01-02 15:04:05.000"))
	hash := sha256.New()
	hash.Write(timeBytes)
	return hex.EncodeToString(hash.Sum(nil))
}

func GetType(f reflect.StructField) (string, error) {
	// available types: text,file,relation,editor,number,bool,email,url,date,select,json
	strings := map[string]string{
		"text":   "text",
		"editor": "text",
	}
	t := f.Tag.Get("type")
	v, found := strings[t]
	if found {
		return v, nil
	}
	if t == "" && f.Type.Kind() == reflect.String {
		return "text", nil
	}
	if t == "" && f.Type.Kind() == reflect.Int {
		return "number", nil
	}
	if t == "" && f.Type.Kind() == reflect.Bool {
		return "bool", nil
	}
	if t == "" && f.Type.Kind() == reflect.Ptr {
		return "relation", nil
	}

	return "", nil
}
