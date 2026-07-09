package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/controller"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/dto"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}

func setupIntegrationRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := repository.NewUserRepositoryGorm(db)
	ctl := controller.NewUserController(repo)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", ctl.ListUser)
		v1.GET("/users/:id", ctl.GetUser)
		v1.POST("/users", ctl.CreateUser)
		v1.DELETE("/users/:id", ctl.DeleteUser)
		v1.PUT("/users/:id", ctl.UpdateUser)
	}

	return router
}

func performRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeResponse(t *testing.T, body []byte) dto.SuccessResponse {
	t.Helper()
	var resp dto.SuccessResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func decodeUser(t *testing.T, body []byte) dto.UserResponse {
	t.Helper()
	resp := decodeResponse(t, body)
	data, _ := json.Marshal(resp.Data)
	var user dto.UserResponse
	require.NoError(t, json.Unmarshal(data, &user))
	return user
}

func decodeUsers(t *testing.T, body []byte) []dto.UserResponse {
	t.Helper()
	resp := decodeResponse(t, body)
	data, _ := json.Marshal(resp.Data)
	var users []dto.UserResponse
	require.NoError(t, json.Unmarshal(data, &users))
	return users
}

func decodeError(t *testing.T, body []byte) dto.ErrorResponse {
	t.Helper()
	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

// --- Integration Tests ---

func TestIntegration_CreateUser(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	rec := performRequest(router, http.MethodPost, "/api/v1/users", body)

	assert.Equal(t, http.StatusCreated, rec.Code)

	user := decodeUser(t, rec.Body.Bytes())
	assert.NotZero(t, user.Id)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)

	var dbUser models.User
	db.First(&dbUser, user.Id)
	assert.NotEqual(t, "secret123", dbUser.Password, "password should be hashed in DB")
	assert.NoError(t, dbUser.VerifyPassword("secret123"))
}

func TestIntegration_CreateUser_InvalidInput(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	tests := []struct {
		name       string
		body       string
		statusCode int
		code       string
	}{
		{"empty body", "", 400, "INVALID_REQUEST"},
		{"invalid json", `{invalid`, 400, "INVALID_REQUEST"},
		{"missing name", `{"email":"a@b.com","password":"pass123"}`, 400, "VALIDATION_ERROR"},
		{"missing email", `{"name":"Alice","password":"pass123"}`, 400, "VALIDATION_ERROR"},
		{"invalid email", `{"name":"Alice","email":"bad","password":"pass123"}`, 400, "VALIDATION_ERROR"},
		{"missing password", `{"name":"Alice","email":"a@b.com"}`, 400, "VALIDATION_ERROR"},
		{"short password", `{"name":"Alice","email":"a@b.com","password":"ab"}`, 400, "VALIDATION_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(router, http.MethodPost, "/api/v1/users", tt.body)
			assert.Equal(t, tt.statusCode, rec.Code)
			errResp := decodeError(t, rec.Body.Bytes())
			assert.Equal(t, tt.code, errResp.Code)
		})
	}

	var count int64
	db.Model(&models.User{}).Count(&count)
	assert.Equal(t, int64(0), count, "no users should be created")
}

func TestIntegration_GetUser_Found(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)
	created := decodeUser(t, createRec.Body.Bytes())

	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	assert.Equal(t, http.StatusOK, getRec.Code)

	got := decodeUser(t, getRec.Body.Bytes())
	assert.Equal(t, created.Id, got.Id)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
}

func TestIntegration_GetUser_NotFound(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	rec := performRequest(router, http.MethodGet, "/api/v1/users/999", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)

	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
}

func TestIntegration_GetUser_InvalidID(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	rec := performRequest(router, http.MethodGet, "/api/v1/users/abc", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestIntegration_ListUsers_Empty(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	rec := performRequest(router, http.MethodGet, "/api/v1/users", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	users := decodeUsers(t, rec.Body.Bytes())
	assert.Empty(t, users)
}

func TestIntegration_ListUsers_Multiple(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Bob","email":"bob@example.com","password":"secret456"}`)
	performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Charlie","email":"charlie@example.com","password":"secret789"}`)

	rec := performRequest(router, http.MethodGet, "/api/v1/users", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	users := decodeUsers(t, rec.Body.Bytes())
	require.Len(t, users, 3)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "Bob", users[1].Name)
	assert.Equal(t, "Charlie", users[2].Name)
}

func TestIntegration_UpdateUser_Found(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Alice","email":"alice@example.com","password":"oldpass123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)
	created := decodeUser(t, createRec.Body.Bytes())

	updateRec := performRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/users/%d", created.Id),
		`{"name":"Alice Updated","email":"alice.new@example.com","password":"newpass123"}`)
	assert.Equal(t, http.StatusOK, updateRec.Code)

	updated := decodeUser(t, updateRec.Body.Bytes())
	assert.Equal(t, created.Id, updated.Id)
	assert.Equal(t, "Alice Updated", updated.Name)
	assert.Equal(t, "alice.new@example.com", updated.Email)

	var dbUser models.User
	db.First(&dbUser, updated.Id)
	assert.Equal(t, "Alice Updated", dbUser.Name)
	assert.Equal(t, "alice.new@example.com", dbUser.Email)
	assert.NoError(t, dbUser.VerifyPassword("newpass123"))
	assert.Error(t, dbUser.VerifyPassword("oldpass123"))
}

func TestIntegration_UpdateUser_NotFound(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	rec := performRequest(router, http.MethodPut, "/api/v1/users/999",
		`{"name":"Ghost","email":"ghost@example.com","password":"pass123"}`)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
}

func TestIntegration_UpdateUser_ValidationErrors(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)
	created := decodeUser(t, createRec.Body.Bytes())

	tests := []struct {
		name string
		body string
	}{
		{"missing name", `{"email":"test@example.com","password":"pass123"}`},
		{"missing email", `{"name":"Alice","password":"pass123"}`},
		{"invalid email", `{"name":"Alice","email":"bad","password":"pass123"}`},
		{"missing password", `{"name":"Alice","email":"a@b.com"}`},
		{"short password", `{"name":"Alice","email":"a@b.com","password":"ab"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/users/%d", created.Id), tt.body)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			errResp := decodeError(t, rec.Body.Bytes())
			assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
		})
	}

	var dbUser models.User
	db.First(&dbUser, created.Id)
	assert.Equal(t, "Alice", dbUser.Name, "user should not be modified")
}

func TestIntegration_DeleteUser_Found(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)
	created := decodeUser(t, createRec.Body.Bytes())

	deleteRec := performRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	assert.Equal(t, http.StatusOK, deleteRec.Code)

	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	assert.Equal(t, http.StatusNotFound, getRec.Code)

	var count int64
	db.Model(&models.User{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestIntegration_DeleteUser_NotFound(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/999", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)

	errResp := decodeError(t, rec.Body.Bytes())
	assert.Equal(t, "USER_NOT_FOUND", errResp.Code)
}

func TestIntegration_DeleteUser_InvalidID(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	rec := performRequest(router, http.MethodDelete, "/api/v1/users/abc", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestIntegration_FullCRUDFlow(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	// Create
	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Flow User","email":"flow@example.com","password":"flowpass123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)
	created := decodeUser(t, createRec.Body.Bytes())
	assert.Equal(t, "Flow User", created.Name)

	// List - should have 1 user
	listRec := performRequest(router, http.MethodGet, "/api/v1/users", "")
	require.Equal(t, http.StatusOK, listRec.Code)
	users := decodeUsers(t, listRec.Body.Bytes())
	require.Len(t, users, 1)

	// Get by ID
	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, getRec.Code)
	got := decodeUser(t, getRec.Body.Bytes())
	assert.Equal(t, "Flow User", got.Name)

	// Update
	updateRec := performRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/users/%d", created.Id),
		`{"name":"Flow Updated","email":"flow.updated@example.com","password":"newflowpass123"}`)
	require.Equal(t, http.StatusOK, updateRec.Code)
	updated := decodeUser(t, updateRec.Body.Bytes())
	assert.Equal(t, "Flow Updated", updated.Name)
	assert.Equal(t, "flow.updated@example.com", updated.Email)

	// Verify update persisted
	getRec2 := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, getRec2.Code)
	got2 := decodeUser(t, getRec2.Body.Bytes())
	assert.Equal(t, "Flow Updated", got2.Name)

	// Delete
	deleteRec := performRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, deleteRec.Code)

	// Verify deleted
	getRec3 := performRequest(router, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", created.Id), "")
	assert.Equal(t, http.StatusNotFound, getRec3.Code)

	// List - should be empty
	listRec2 := performRequest(router, http.MethodGet, "/api/v1/users", "")
	require.Equal(t, http.StatusOK, listRec2.Code)
	users2 := decodeUsers(t, listRec2.Body.Bytes())
	assert.Empty(t, users2)
}

func TestIntegration_PasswordNeverExposed(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	createRec := performRequest(router, http.MethodPost, "/api/v1/users",
		`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	require.Equal(t, http.StatusCreated, createRec.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &body))
	data, _ := json.Marshal(body["data"])
	var userData map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &userData))
	_, hasPassword := userData["password"]
	assert.False(t, hasPassword, "password should not be in JSON response")

	getRec := performRequest(router, http.MethodGet, "/api/v1/users/1", "")
	require.Equal(t, http.StatusOK, getRec.Code)

	var getBody map[string]interface{}
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &getBody))
	getData, _ := json.Marshal(getBody["data"])
	var getUser map[string]interface{}
	require.NoError(t, json.Unmarshal(getData, &getUser))
	_, hasGetPassword := getUser["password"]
	assert.False(t, hasGetPassword, "password should not be in GET response")
}

func TestIntegration_ConcurrentCreate(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db)

	const numUsers = 5
	done := make(chan bool, numUsers)

	for i := 0; i < numUsers; i++ {
		go func(idx int) {
			body := fmt.Sprintf(`{"name":"User%d","email":"user%d@example.com","password":"pass%d123"}`, idx, idx, idx)
			rec := performRequest(router, http.MethodPost, "/api/v1/users", body)
			done <- rec.Code == http.StatusCreated
		}(i)
	}

	for i := 0; i < numUsers; i++ {
		<-done
	}

	var count int64
	db.Model(&models.User{}).Count(&count)
	assert.True(t, count > 0, "at least some users should be created")
}
