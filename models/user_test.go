package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserTableName(t *testing.T) {
	user := User{}
	if got := user.TableName(); got != "table_of_user" {
		t.Errorf("TableName() = %q, want %q", got, "table_of_user")
	}
}

func TestUser_StructFields(t *testing.T) {
	user := User{
		Id:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "hashed",
	}

	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, "hashed", user.Password)
}

func TestHashPassword(t *testing.T) {
	user := &User{}
	err := user.HashPassword("mypassword")

	require.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "mypassword", user.Password)
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	user1 := &User{}
	user1.HashPassword("samepass")

	user2 := &User{}
	user2.HashPassword("samepass")

	assert.NotEqual(t, user1.Password, user2.Password)
}

func TestVerifyPassword_Correct(t *testing.T) {
	user := &User{}
	plainPassword := "secret123"
	err := user.HashPassword(plainPassword)
	require.NoError(t, err)

	err = user.VerifyPassword(plainPassword)
	assert.NoError(t, err)
}

func TestVerifyPassword_Wrong(t *testing.T) {
	user := &User{}
	err := user.HashPassword("correctpass")
	require.NoError(t, err)

	err = user.VerifyPassword("wrongpass")
	assert.Error(t, err)
}

func TestVerifyPassword_EmptyPassword(t *testing.T) {
	user := &User{}
	err := user.HashPassword("password")
	require.NoError(t, err)

	err = user.VerifyPassword("")
	assert.Error(t, err)
}
