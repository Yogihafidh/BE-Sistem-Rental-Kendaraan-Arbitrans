package utils

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "Secret123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hash == password {
		t.Fatal("hash should not match plain password")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("expected password to match hash")
	}

	if CheckPasswordHash("WrongPassword", hash) {
		t.Fatal("expected wrong password to not match hash")
	}
}
