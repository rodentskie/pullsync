package mapstruct

import (
	"encoding/json"
	"reflect"
)

func MapToStruct(data map[string]interface{}, targetStruct interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, targetStruct); err != nil {
		return err
	}

	return nil
}

func StructToMapInterface(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		key := field.Tag.Get("json")
		value := val.Field(i).Interface()
		result[key] = value
	}

	return result
}
