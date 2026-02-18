package auth_test

import (
	"testing"

	"github.com/andriyzaec/chirpy/internal/auth"
)

func TestHashAndCheckPassword_Success(t *testing.T) {
	password := "mySecret123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error while hashing password: %v", err)
	}

	ok, err := auth.CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("unexpected error while checking password: %v", err)
	}

	if !ok {
		t.Fatal("expected password to match hash, but it did not")
	}
}

func TestHashAndCheckPassword_WrongPassword(t *testing.T) {
	password := "mySecret123"
	wrongPassword := "wrongPassword"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error while hashing password: %v", err)
	}

	ok, err := auth.CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("unexpected error while checking password: %v", err)
	}

	if ok {
		t.Fatal("expected password NOT to match hash, but it did")
	}
}
