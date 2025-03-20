package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// Test JWT creation
func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	if token == "" {
		t.Errorf("Generated token is empty")
	}
}

// Test JWT validation
func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedUserID)
	}
}

// Test invalid JWT validation
func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("invalid.token.string", "testsecret")
	if err == nil {
		t.Errorf("Expected an error for invalid token, got nil")
	}
}
