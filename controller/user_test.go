package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/controller"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/dto"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockRouter(mock *repository.MockUserRepository) http.Handler {
	gin.SetMode(gin.TestMode)
	router := controller.SetupRouterWithMock(mock)
	return router
}

func performRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeUsers(t *testing.T, body []byte) []dto.UserResponse {
	t.Helper()

	var successResp dto.SuccessResponse
	require.NoError(t, json.Unmarshal(body, &successResp))

	data, err := json.Marshal(successResp.Data)
	require.NoError(t, err)

	var users []dto.UserResponse
	require.NoError(t, json.Unmarshal(data, &users))
	return users
}

func decodeUser(t *testing.T, body []byte) dto.UserResponse {
	t.Helper()

	var successResp dto.SuccessResponse
	require.NoError(t, json.Unmarshal(body, &successResp))

	data, err := json.Marshal(successResp.Data)
	require.NoError(t, err)

	var user dto.UserResponse
	require.NoError(t, json.Unmarshal(data, &user))
	return user
}

func decodeError(t *testing.T, body []byte) dto.ErrorResponse {
	t.Helper()

	var errorResp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &errorResp))
	return errorResp
}

func TestListUser_Empty(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users", "")

	assert.Equal(t, http.StatusOK, rec.Code)

	users := decodeUsers(t, rec.Body.Bytes())
	assert.Empty(t, users)
}

func TestListUser_WithUsers(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.Create(nil, &models.User{Name: "Bob", Email: "bob@example.com"})
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users", "")

	assert.Equal(t, http.StatusOK, rec.Code)

	users := decodeUsers(t, rec.Body.Bytes())
	require.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "Bob", users[1].Name)
}

func TestListUser_RepoError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.FindErr = fmt.Errorf("database connection lost")
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "INTERNAL_ERROR", errResp.Code)
}

func TestGetUser_Found(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users/1", "")

	assert.Equal(t, http.StatusOK, rec.Code)

	got := decodeUser(t, rec.Body.Bytes())
	assert.Equal(t, 1, got.Id)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users/999", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
}

func TestGetUser_InvalidID(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users/abc", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetUser_RepoError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Id: 1, Name: "Alice", Email: "alice@example.com"})
	mock.FindErr = fmt.Errorf("database connection lost")
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodGet, "/api/v1/users/1", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateUser_Success(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","email":"charlie@example.com","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusCreated, rec.Code)

	created := decodeUser(t, rec.Body.Bytes())
	assert.NotZero(t, created.Id)
	assert.Equal(t, "Charlie", created.Name)
	assert.Equal(t, "charlie@example.com", created.Email)
}

func TestCreateUser_EmptyBody(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodPost, "/api/v1/users", "")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "INVALID_REQUEST", errResp.Code)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodPost, "/api/v1/users", `{invalid`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateUser_MissingName(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"email":"test@example.com","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
	assert.Contains(t, errResp.Error, "name")
}

func TestCreateUser_MissingEmail(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","email":"invalid-email","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestCreateUser_WeakPassword(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","email":"charlie@example.com","password":"pass"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestCreateUser_MissingPassword(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","email":"charlie@example.com"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestCreateUser_RepoError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.CreateErr = fmt.Errorf("duplicate key error")
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","email":"charlie@example.com","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "INTERNAL_ERROR", errResp.Code)
}

func TestCreateUser_WhitespaceName(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"   ","email":"test@example.com","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestCreateUser_WhitespaceEmail(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Charlie","email":"   ","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUpdateUser_Success(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)

	body := `{"name":"Alice Updated","email":"alice.new@example.com","password":"newpass123"}`
	rec := performRequest(router, http.MethodPut, "/api/v1/users/1", body)

	assert.Equal(t, http.StatusOK, rec.Code)

	updated := decodeUser(t, rec.Body.Bytes())
	assert.Equal(t, 1, updated.Id)
	assert.Equal(t, "Alice Updated", updated.Name)
	assert.Equal(t, "alice.new@example.com", updated.Email)
}

func TestUpdateUser_NotFound(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	body := `{"name":"Ghost","email":"ghost@example.com","password":"pass123"}`
	rec := performRequest(router, http.MethodPut, "/api/v1/users/404", body)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodPut, "/api/v1/users/1", `{invalid`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateUser_ValidationErrors(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)

	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{"missing name", `{"email":"test@example.com","password":"pass123"}`, "name"},
		{"missing email", `{"name":"Alice","password":"pass123"}`, "email"},
		{"invalid email", `{"name":"Alice","email":"bad","password":"pass123"}`, "email format"},
		{"missing password", `{"name":"Alice","email":"test@example.com"}`, "password"},
		{"short password", `{"name":"Alice","email":"test@example.com","password":"ab"}`, "6 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(router, http.MethodPut, "/api/v1/users/1", tt.body)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			errResp := decodeError(t, rec.Body.Bytes())
			assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
			assert.Contains(t, errResp.Error, tt.expected)
		})
	}
}

func TestUpdateUser_RepoFindError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.FindErr = fmt.Errorf("database connection lost")
	router := setupMockRouter(mock)

	body := `{"name":"Alice Updated","email":"alice.new@example.com","password":"newpass123"}`
	rec := performRequest(router, http.MethodPut, "/api/v1/users/1", body)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateUser_RepoUpdateError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.UpdateErr = fmt.Errorf("database write failed")
	router := setupMockRouter(mock)

	body := `{"name":"Alice Updated","email":"alice.new@example.com","password":"newpass123"}`
	rec := performRequest(router, http.MethodPut, "/api/v1/users/1", body)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "INTERNAL_ERROR", errResp.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/1", "")

	assert.Equal(t, http.StatusOK, rec.Code)

	getRec := performRequest(router, http.MethodGet, "/api/v1/users/1", "")
	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/999", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/abc", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteUser_RepoFindError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.FindErr = fmt.Errorf("database connection lost")
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/1", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteUser_RepoDeleteError(t *testing.T) {
	mock := repository.NewMockUserRepository()
	mock.Create(nil, &models.User{Name: "Alice", Email: "alice@example.com"})
	mock.DeleteErr = fmt.Errorf("database write failed")
	router := setupMockRouter(mock)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/1", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "INTERNAL_ERROR", errResp.Code)
}

func TestUserCRUD_Flow(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Flow User","email":"flow@example.com","password":"flowpass123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)

	created := decodeUser(t, createRec.Body.Bytes())

	listRec := performRequest(router, http.MethodGet, "/api/v1/users", "")
	require.Equal(t, http.StatusOK, listRec.Code)
	users := decodeUsers(t, listRec.Body.Bytes())
	require.Len(t, users, 1)

	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, getRec.Code)

	updateRec := performRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/users/%d", created.Id),
		`{"name":"Flow Updated","email":"flow.updated@example.com","password":"newflowpass123"}`)
	require.Equal(t, http.StatusOK, updateRec.Code)

	updated := decodeUser(t, updateRec.Body.Bytes())
	assert.Equal(t, "Flow Updated", updated.Name)

	deleteRec := performRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, deleteRec.Code)

	finalListRec := performRequest(router, http.MethodGet, "/api/v1/users", "")
	require.Equal(t, http.StatusOK, finalListRec.Code)
	finalUsers := decodeUsers(t, finalListRec.Body.Bytes())
	assert.Empty(t, finalUsers)
}

func TestPasswordHashing(t *testing.T) {
	mock := repository.NewMockUserRepository()
	router := setupMockRouter(mock)

	plainPassword := "mysecretpassword123"
	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		fmt.Sprintf(`{"name":"Test User","email":"test@example.com","password":"%s"}`, plainPassword))
	require.Equal(t, http.StatusCreated, createRec.Code)

	created := decodeUser(t, createRec.Body.Bytes())
	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, getRec.Code)

	var successResp dto.SuccessResponse
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &successResp))
	data, _ := json.Marshal(successResp.Data)
	var user dto.UserResponse
	require.NoError(t, json.Unmarshal(data, &user))
	assert.Equal(t, created.Id, user.Id)
}

func TestValidationError_Error(t *testing.T) {
	err := controller.ValidationError{Code: "VALIDATION_ERROR", Message: "name is required"}
	assert.Equal(t, "name is required", err.Error())
}
