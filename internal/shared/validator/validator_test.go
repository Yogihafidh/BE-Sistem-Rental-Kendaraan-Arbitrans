package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type sampleRequest struct {
	Email    string `validate:"required,email"`
	FullName string `validate:"required,min=3"`
	Age      int    `validate:"gte=18"`
}

func TestFormatValidationErrors(t *testing.T) {
	validate := validator.New()

	req := sampleRequest{
		Email:    "not-an-email",
		FullName: "ab", // min 3
		Age:      15,   // gte 18
	}

	err := validate.Struct(req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	fieldErrors := FormatValidationErrors(err)
	if len(fieldErrors) != 3 {
		t.Fatalf("expected 3 field errors, got %d", len(fieldErrors))
	}

	fieldMap := make(map[string]FieldError)
	for _, fe := range fieldErrors {
		fieldMap[fe.Field] = fe
	}

	if _, ok := fieldMap["email"]; !ok {
		t.Error("expected error for field 'email'")
	}
	if _, ok := fieldMap["full_name"]; !ok {
		t.Error("expected error for field 'full_name'")
	}
	if _, ok := fieldMap["age"]; !ok {
		t.Error("expected error for field 'age'")
	}
}
