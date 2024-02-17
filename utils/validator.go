package utils

import (
	"gopkg.in/go-playground/validator.v8"
)

type Validate interface {
	Validate() error
}

func ValidateStruct(obj Validate) error {
	config := &validator.Config{TagName: "validate"}

	validate := validator.New(config)

	err := validate.Struct(obj)
	if err != nil {
		return err
	}

	return nil
}
