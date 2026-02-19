package auth_test

import (
	"testing"
	"time"

	"github.com/andriyzaec/chirpy/internal/auth"
	"github.com/google/uuid"
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

func TestJWT_Success(t *testing.T) {
	supersecret := "supersecret"
	uuid, _ := uuid.NewUUID()
	token, err := auth.MakeJWT(uuid, supersecret, 60*time.Second)
	if err != nil {
		t.Fatalf("unexpected error while creating token: %v", err)
	}

	claimsID, err := auth.ValidateJWT(token, supersecret)
	if err != nil {
		t.Fatalf("unexpected error while validating JWT: %v", err)
	}

	if claimsID != uuid {
		t.Fatal("expected claimsId = uuid, not fulfilled")
	}
}

func TestJWT_BadSecret(t *testing.T) {
	supersecret := "supersecret"
	uuid, _ := uuid.NewUUID()
	token, err := auth.MakeJWT(uuid, supersecret, 60*time.Second)
	if err != nil {
		t.Fatalf("unexpected error while creating token: %v", err)
	}

	badSecret := "bad_secret"
	_, err = auth.ValidateJWT(token, badSecret)
	if err == nil {
		t.Fatal("No error when bad secret is using, should be error")
	}
}
