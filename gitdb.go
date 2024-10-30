package gitdb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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

func Unmarshal(data []byte, v any) error {
	docs, _, err := DecodeDocs(data)
	if err != nil {
		return err
	}

	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Pointer {
		return fmt.Errorf("object to unmarshal must be a pointer")
	}
	if len(docs) == 0 {
		return fmt.Errorf("no documents to unmarshal")
	}
	if len(docs) == 1 {
		bytes, _ := json.Marshal(docs[0])
		return json.Unmarshal(bytes, v)
	}
	if len(docs) > 1 {
		bytes, _ := json.Marshal(docs)
		return json.Unmarshal(bytes, v)
	}
	return nil
}

func DecodeDocs(data []byte) ([]map[string]interface{}, map[string]bool, error) {
	docs := make([]map[string]interface{}, 0)
	lines := strings.Split(string(data), "\n")
	in_doc := false
	doc := make(map[string]interface{}, 0)
	meta := make(map[string]bool, 0)
	buffer := []string{}
	key := ""
	is_meta := false
	for _, line := range lines {
		if strings.HasPrefix(line, "//gitdb:doc:begin") {
			in_doc = true
			continue
		}
		if strings.HasPrefix(line, "//gitdb:doc:end") {
			if len(buffer) > 0 {
				doc[key] = strings.Join(buffer, "\n")
			}
			docs = append(docs, doc)
			doc = make(map[string]interface{}, 0)
			in_doc = false
			continue
		}
		if in_doc {
			if strings.HasPrefix(line, "//gitdb:field") {
				if len(buffer) > 0 {
					doc[key] = strings.Join(buffer, "\n")
				}
				key = strings.Split(line, ":")[2]
				buffer = make([]string, 0)
				is_meta = false
				meta[key] = false
				continue
			}
			if strings.HasPrefix(line, "//gitdb:meta") {
				is_meta = true
				if len(buffer) > 0 {
					doc[key] = strings.Join(buffer, "\n")
				}
				key = strings.Split(line, ":")[2]
				buffer = make([]string, 0)
				is_meta = true
				meta[key] = true
				continue
			}
			if is_meta {
				val := line[2:]
				doc[key] = val
				buffer = make([]string, 0)
			} else {
				buffer = append(buffer, line)
			}
		}
	}

	return docs, meta, nil
}
