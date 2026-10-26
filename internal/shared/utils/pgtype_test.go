package utils

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestUUIDConverters(t *testing.T) {
	rawUUID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	// 1. RequiredUUID valid
	u, err := RequiredUUID(rawUUID)
	if err != nil {
		t.Fatalf("unexpected error scanning valid UUID: %v", err)
	}
	if !u.Valid {
		t.Fatal("expected UUID to be valid")
	}

	// 2. UUIDToString
	str := UUIDToString(u)
	if str != rawUUID {
		t.Errorf("expected '%s', got '%s'", rawUUID, str)
	}

	// 3. UUIDToString invalid
	invalidU := pgtype.UUID{Valid: false}
	if UUIDToString(invalidU) != "" {
		t.Errorf("expected empty string for invalid UUID, got '%s'", UUIDToString(invalidU))
	}

	// 4. RequiredUUID invalid
	_, err = RequiredUUID("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID string")
	}

	// 5. OptionalUUID empty
	optU, err := OptionalUUID("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if optU.Valid {
		t.Fatal("expected invalid UUID for empty string")
	}

	// 6. OptionalUUID populated
	optU2, err := OptionalUUID(rawUUID)
	if err != nil || !optU2.Valid {
		t.Fatalf("expected valid UUID for populated string: %v", err)
	}
}

func TestTextConverters(t *testing.T) {
	// 1. Valid text
	txt := pgtype.Text{String: "hello world", Valid: true}
	ptr := Text(txt)
	if ptr == nil || *ptr != "hello world" {
		t.Errorf("expected 'hello world', got %v", ptr)
	}

	// 2. Invalid text
	invTxt := pgtype.Text{Valid: false}
	if Text(invTxt) != nil {
		t.Errorf("expected nil for invalid pgtype.Text")
	}

	// 3. OptionalText nil or empty
	if OptionalText(nil).Valid {
		t.Error("expected invalid for nil string pointer")
	}
	emptyStr := ""
	if OptionalText(&emptyStr).Valid {
		t.Error("expected invalid for empty string pointer")
	}

	// 4. OptionalText valid
	val := "arbitrans"
	optRes := OptionalText(&val)
	if !optRes.Valid || optRes.String != "arbitrans" {
		t.Errorf("expected 'arbitrans', got %v", optRes)
	}
}

func TestInt4Converters(t *testing.T) {
	// 1. Valid Int4
	i := pgtype.Int4{Int32: 42, Valid: true}
	ptr := Int4(i)
	if ptr == nil || *ptr != 42 {
		t.Errorf("expected 42, got %v", ptr)
	}

	// 2. Invalid Int4
	inv := pgtype.Int4{Valid: false}
	if Int4(inv) != nil {
		t.Errorf("expected nil for invalid Int4")
	}

	// 3. OptionalInt4
	if OptionalInt4(nil).Valid {
		t.Error("expected invalid for nil int32 pointer")
	}
	val := int32(100)
	opt := OptionalInt4(&val)
	if !opt.Valid || opt.Int32 != 100 {
		t.Errorf("expected 100, got %v", opt)
	}
}

func TestDateConverters(t *testing.T) {
	dateStr := "2026-10-26"

	// 1. RequiredDate
	d, err := RequiredDate(dateStr)
	if err != nil {
		t.Fatalf("unexpected error parsing date: %v", err)
	}
	if !d.Valid {
		t.Fatal("expected date to be valid")
	}

	// 2. Date to string
	strPtr := Date(d)
	if strPtr == nil || *strPtr != dateStr {
		t.Errorf("expected '%s', got %v", dateStr, strPtr)
	}

	// 3. Invalid Date
	invD := pgtype.Date{Valid: false}
	if Date(invD) != nil {
		t.Error("expected nil for invalid date")
	}

	// 4. OptionalDate
	optD, err := OptionalDate(nil)
	if err != nil || optD.Valid {
		t.Error("expected invalid for nil date string")
	}
	empty := ""
	optD2, err := OptionalDate(&empty)
	if err != nil || optD2.Valid {
		t.Error("expected invalid for empty date string")
	}
	optD3, err := OptionalDate(&dateStr)
	if err != nil || !optD3.Valid {
		t.Error("expected valid for date string")
	}
}

func TestTimestamptzConverter(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	ts := pgtype.Timestamptz{Time: now, Valid: true}

	res := Timestamptz(ts)
	if res == nil || !res.Equal(now) {
		t.Errorf("expected %v, got %v", now, res)
	}

	inv := pgtype.Timestamptz{Valid: false}
	if Timestamptz(inv) != nil {
		t.Error("expected nil for invalid Timestamptz")
	}
}

func TestNumericConverters(t *testing.T) {
	// 1. NumericFromString valid
	num, err := NumericFromString("1500000.50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !num.Valid {
		t.Fatal("expected numeric to be valid")
	}

	// 2. Numeric to string
	str := Numeric(num)
	if str != "1500000.5" {
		t.Errorf("expected '1500000.5', got '%s'", str)
	}

	// 3. Invalid Numeric
	invNum := pgtype.Numeric{Valid: false}
	if Numeric(invNum) != "0" {
		t.Errorf("expected '0' for invalid Numeric, got '%s'", Numeric(invNum))
	}
}
