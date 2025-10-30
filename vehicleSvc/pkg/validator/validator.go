package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(obj interface{}) error {
	err := validate.Struct(obj)
	return customError(err)
}

func customError(err error) error {
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				return NewValidationError(fmt.Sprintf("%s is required",
					err.Field()))
			case "email":
				return NewValidationError(fmt.Sprintf("%s is not valid email",
					err.Field()))
			case "unique":
				return NewValidationError(fmt.Sprintf("%s must unique value",
					err.Field()))
			case "notEmpty":
				return NewValidationError(fmt.Sprintf("%s can not be empty",
					err.Field()))
			case "max":
				return NewValidationError(fmt.Sprintf("%s value must be lower than %s", err.Field(), err.Param()))
			case "min":
				return NewValidationError(fmt.Sprintf("%s value must be grather than %s", err.Field(), err.Param()))
			default:
				return NewValidationError(fmt.Sprintf("%s validation error on %s tag", err.Field(), err.ActualTag()))
			}
		}
	}
	return nil
}

func NewValidationError(msg string) ValidationErrors {
	return ValidationErrors{errors.New(msg)}
}

type ValidationErrors struct {
	err error
}

func (v ValidationErrors) Error() string {
	return v.err.Error()
}
