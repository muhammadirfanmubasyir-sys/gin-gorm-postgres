package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSuccessResponse(t *testing.T) {
	data := UserResponse{Id: 1, Name: "Alice", Email: "alice@example.com"}
	resp := NewSuccessResponse(data)

	assert.Equal(t, data, resp.Data)
	assert.WithinDuration(t, time.Now(), resp.Timestamp, time.Second)
}

func TestNewSuccessResponse_NilData(t *testing.T) {
	resp := NewSuccessResponse(nil)

	assert.Nil(t, resp.Data)
	assert.WithinDuration(t, time.Now(), resp.Timestamp, time.Second)
}

func TestNewSuccessResponse_SliceData(t *testing.T) {
	users := []UserResponse{
		{Id: 1, Name: "Alice", Email: "alice@example.com"},
		{Id: 2, Name: "Bob", Email: "bob@example.com"},
	}
	resp := NewSuccessResponse(users)

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded SuccessResponse
	require.NoError(t, json.Unmarshal(data, &decoded))

	marshaled, _ := json.Marshal(decoded.Data)
	var decodedUsers []UserResponse
	require.NoError(t, json.Unmarshal(marshaled, &decodedUsers))
	require.Len(t, decodedUsers, 2)
	assert.Equal(t, "Alice", decodedUsers[0].Name)
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse("not found", "NOT_FOUND")

	assert.Equal(t, "not found", resp.Error)
	assert.Equal(t, "NOT_FOUND", resp.Code)
	assert.WithinDuration(t, time.Now(), resp.Timestamp, time.Second)
}

func TestNewErrorResponse_EmptyStrings(t *testing.T) {
	resp := NewErrorResponse("", "")

	assert.Equal(t, "", resp.Error)
	assert.Equal(t, "", resp.Code)
	assert.WithinDuration(t, time.Now(), resp.Timestamp, time.Second)
}

func TestErrorResponse_JSON(t *testing.T) {
	resp := NewErrorResponse("validation failed", "VALIDATION_ERROR")

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded ErrorResponse
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, "validation failed", decoded.Error)
	assert.Equal(t, "VALIDATION_ERROR", decoded.Code)
	assert.False(t, decoded.Timestamp.IsZero())
}

func TestSuccessResponse_JSON(t *testing.T) {
	resp := NewSuccessResponse(UserResponse{Id: 1, Name: "Alice", Email: "alice@example.com"})

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))

	_, hasData := decoded["data"]
	_, hasTimestamp := decoded["timestamp"]
	assert.True(t, hasData)
	assert.True(t, hasTimestamp)
}

func TestUserResponse_JSON(t *testing.T) {
	user := UserResponse{Id: 42, Name: "Charlie", Email: "charlie@test.com"}

	data, err := json.Marshal(user)
	require.NoError(t, err)

	var decoded UserResponse
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, 42, decoded.Id)
	assert.Equal(t, "Charlie", decoded.Name)
	assert.Equal(t, "charlie@test.com", decoded.Email)
}

func TestCreateUserRequest_JSON(t *testing.T) {
	jsonStr := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	var req CreateUserRequest
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &req))

	assert.Equal(t, "Alice", req.Name)
	assert.Equal(t, "alice@example.com", req.Email)
	assert.Equal(t, "secret123", req.Password)
}

func TestUpdateUserRequest_JSON(t *testing.T) {
	jsonStr := `{"name":"Updated","email":"updated@example.com","password":"newpass123"}`
	var req UpdateUserRequest
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &req))

	assert.Equal(t, "Updated", req.Name)
	assert.Equal(t, "updated@example.com", req.Email)
	assert.Equal(t, "newpass123", req.Password)
}
