package utils

import (
	"fmt"
	"reflect"
)

func CheckAndUpdate[T any, G any](payload T, entity *G) error {
	payloadValue := reflect.ValueOf(payload)
	entityValue := reflect.ValueOf(entity).Elem()

	if payloadValue.Kind() != reflect.Struct || entityValue.Kind() != reflect.Struct {
		return fmt.Errorf("invalid data types for update")
	}

	for i := 0; i < payloadValue.NumField(); i++ {
		fieldValue := payloadValue.Field(i)
		fieldName := payloadValue.Type().Field(i).Name

		if fieldValue.IsValid() && !fieldValue.IsZero() {
			entityField := entityValue.FieldByName(fieldName)
			if entityField.IsValid() && entityField.CanSet() {
				entityField.Set(fieldValue)
			}
		}
	}
	return nil
}
