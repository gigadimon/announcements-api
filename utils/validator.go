package utils

import (
	"errors"
	"reflect"

	"gopkg.in/go-playground/validator.v8"
)

func ValidateStruct(obj interface{}) error {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr && objType.Elem().Kind() != reflect.Struct {
		return errors.New("passed data isn't struct type")
	}

	config := &validator.Config{TagName: "validate"}

	validate := validator.New(config)

	err := validate.Struct(obj)
	if err != nil {
		return err
	}

	return nil
}
