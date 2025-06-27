package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	syncOnce sync.Once
)

func init() {
	syncOnce.Do(func() {
		validate = validator.New()
	})
}

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return &ValidationError{
				Field:   getFieldName(validationErrors[0]),
				Message: getValidationMessage(validationErrors[0]),
			}
		}
		return err
	}
	return nil
}

func getFieldName(fieldError validator.FieldError) string {
	return strings.ToLower(fieldError.Field())
}

func getValidationMessage(fieldError validator.FieldError) string {
	field := strings.ToLower(fieldError.Field())
	tag := fieldError.Tag()
	param := fieldError.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "oneof":
		if field == "status" {
			return "status must be 0 (incomplete) or 1 (complete)"
		}
		return fmt.Sprintf("%s must be one of: %s", field, param)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
