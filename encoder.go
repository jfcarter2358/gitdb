package gitdb

import (
	"fmt"
	"reflect"
)

func Marshal(v any) ([]byte, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	out := "//gitdb:doc:begin\n"

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i).Interface()

		field, found := typ.FieldByName(val.Type().Field(i).Name)
		if !found {
			continue
		}

		keyType := "field"
		valPrefix := ""
		meta := field.Tag.Get("gitdb_meta")
		if meta == "true" {
			keyType = "meta"
			valPrefix = "//"
		}
		key := field.Tag.Get("json")
		if key != "" {
			out += fmt.Sprintf("//gitdb:%s:%s\n", keyType, key)
			out += fmt.Sprintf("%s%s\n", valPrefix, value)
		}
	}

	out += "//gitdb:doc:end\n"

	return []byte(out), nil
}
