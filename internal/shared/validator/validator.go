package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldError represents a standardized validation error for a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}

// FormatValidationErrors converts go-playground validator errors into structured, user-friendly field errors.
func FormatValidationErrors(err error) []FieldError {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		result := make([]FieldError, 0, len(validationErrors))
		for _, fe := range validationErrors {
			field := toSnakeCase(fe.Field())
			msg := getErrorMessage(field, fe)
			result = append(result, FieldError{
				Field:   field,
				Message: msg,
				Tag:     fe.Tag(),
			})
		}
		return result
	}

	// Fallback for non-validator errors (e.g., JSON syntax error)
	return []FieldError{
		{
			Field:   "body",
			Message: err.Error(),
		},
	}
}

// getErrorMessage constructs a human-readable validation error message.
func getErrorMessage(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters/items", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters/items", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s failed validation on '%s'", field, fe.Tag())
	}
}

// toSnakeCase converts CamelCase field names to snake_case.
func toSnakeCase(str string) string {
	var sb strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			sb.WriteRune('_')
		}
		sb.WriteRune(r)
	}
	return strings.ToLower(sb.String())
}
