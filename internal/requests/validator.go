package requests

import (
	"fmt"
	"strings"
	"sync"
	"tasks-service-demo/internal/errors"
	apperrors "tasks-service-demo/internal/errors"

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

// ValidateStruct validates a struct using go-playground/validator and returns an AppError if validation fails.
// It converts validation errors to structured AppError responses with appropriate error codes.
func ValidateStruct(s interface{}) *apperrors.AppError {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			message := getValidationMessage(validationErrors[0])
			code := getValidationErrorCode(validationErrors[0])
			return apperrors.NewValidationError(code, message)
		}
		return apperrors.ErrTaskInvalidInput.WithCause(err)
	}
	return nil
}

// getValidationMessage converts a validator field error to a human-readable message.
// It handles different validation tags like required, min, max, and oneof.
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

// getValidationErrorCode maps validator field errors to application error codes.
// It provides specific error codes for different validation failures.
func getValidationErrorCode(fieldError validator.FieldError) int {
	field := strings.ToLower(fieldError.Field())
	tag := fieldError.Tag()

	switch {
	case tag == "required" && field == "name":
		return errors.ErrCodeTaskNameRequired
	case tag == "max" && field == "name":
		return errors.ErrCodeTaskNameTooLong
	case field == "status":
		return errors.ErrCodeTaskInvalidStatus
	default:
		return errors.ErrCodeTaskInvalidInput
	}
}
