package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_AllFieldsEmpty(t *testing.T) {
	err := validate("", "", "")
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Contains(t, err.Error(), "name")
}

func TestValidate_OnlyNameProvided(t *testing.T) {
	err := validate("Alice", "", "")
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Contains(t, err.Error(), "email")
}

func TestValidate_InvalidEmailFormat(t *testing.T) {
	err := validate("Alice", "not-an-email", "pass123")
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Contains(t, err.Error(), "email format")
}

func TestValidate_EmptyPassword(t *testing.T) {
	err := validate("Alice", "alice@example.com", "")
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Contains(t, err.Error(), "password")
}

func TestValidate_ShortPassword(t *testing.T) {
	err := validate("Alice", "alice@example.com", "ab")
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Contains(t, err.Error(), "6 characters")
}

func TestValidate_Success(t *testing.T) {
	err := validate("Alice", "alice@example.com", "pass123")
	assert.Equal(t, "", err.Code)
}
