package utils

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Function to convert UUID to string
func UUIDToString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	return value.String()
}

// Function to convert Text to string
func Text(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

// Function to convert Int4 to int32
func Int4(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}

// Function to convert Date to string
func Date(value pgtype.Date) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format("2006-01-02")
	return &formatted
}

// Function to convert Timestamptz to time.Time
func Timestamptz(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

// Function to convert Numeric to string
func Numeric(value pgtype.Numeric) string {
	if !value.Valid {
		return "0"
	}
	number, err := value.Float64Value()
	if err == nil && number.Valid {
		return strconv.FormatFloat(number.Float64, 'f', -1, 64)
	}
	return "0"
}

// Function to convert string to UUID
func RequiredUUID(value string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(value)
	return uuid, err
}

// Function to convert string to UUID (optional)
func OptionalUUID(value string) (pgtype.UUID, error) {
	if value == "" {
		return pgtype.UUID{}, nil
	}
	return RequiredUUID(value)
}

// Function to convert string to Text (optional)
func OptionalText(value *string) pgtype.Text {
	if value == nil || *value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

// Function to convert string to Date (optional)
func OptionalDate(value *string) (pgtype.Date, error) {
	if value == nil || *value == "" {
		return pgtype.Date{}, nil
	}
	date, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return pgtype.Date{}, err
	}
	return pgtype.Date{Time: date, Valid: true}, nil
}

// Function to convert string to Date (required)
func RequiredDate(value string) (pgtype.Date, error) {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, err
	}
	return pgtype.Date{Time: date, Valid: true}, nil
}

// Function to convert string to Int4 (optional)
func OptionalInt4(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}

// Function to convert string to Numeric
func NumericFromString(value string) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	err := numeric.Scan(value)
	return numeric, err
}
