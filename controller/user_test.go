package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func performRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeUsers(t *testing.T, body []byte) []models.User {
	t.Helper()

	var users []models.User
	require.NoError(t, json.Unmarshal(body, &users))
	return users
}

func decodeUser(t *testing.T, body []byte) models.User {
	t.Helper()

	var user models.User
	require.NoError(t, json.Unmarshal(body, &user))
	return user
}

func TestListUser_Empty(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	rec := performRequest(router, http.MethodGet, "/users", "")

	assert.Equal(t, http.StatusOK, rec.Code)

	users := decodeUsers(t, rec.Body.Bytes())
	assert.Empty(t, users)
}

func TestListUser_WithUsers(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	testutil.SeedUser(t, db, "Alice", "alice@example.com", "secret1")
	testutil.SeedUser(t, db, "Bob", "bob@example.com", "secret2")

	rec := performRequest(router, http.MethodGet, "/users", "")

	assert.Equal(t, http.StatusOK, rec.Code)

	users := decodeUsers(t, rec.Body.Bytes())
	require.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "Bob", users[1].Name)
}

func TestGetUser_Found(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	user := testutil.SeedUser(t, db, "Alice", "alice@example.com", "secret")

	rec := performRequest(router, http.MethodGet, fmt.Sprintf("/users/%d", user.Id), "")

	assert.Equal(t, http.StatusOK, rec.Code)

	got := decodeUser(t, rec.Body.Bytes())
	assert.Equal(t, user.Id, got.Id)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
	assert.Equal(t, "secret", got.Password)
}

func TestGetUser_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	rec := performRequest(router, http.MethodGet, "/users/999", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `"Cannot find user with Id = 999"`, rec.Body.String())
}

func TestGetUser_InvalidID(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	rec := performRequest(router, http.MethodGet, "/users/abc", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `"Cannot find user with Id = abc"`, rec.Body.String())
}

func TestCreateUser_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	body := `{"name":"Charlie","email":"charlie@example.com","password":"pass123"}`
	rec := performRequest(router, http.MethodPost, "/users", body)

	assert.Equal(t, http.StatusOK, rec.Code)

	created := decodeUser(t, rec.Body.Bytes())
	assert.NotZero(t, created.Id)
	assert.Equal(t, "Charlie", created.Name)
	assert.Equal(t, "charlie@example.com", created.Email)
	assert.Equal(t, "pass123", created.Password)
}

func TestCreateUser_EmptyBody(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	rec := performRequest(router, http.MethodPost, "/users", "")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	rec := performRequest(router, http.MethodPost, "/users", `{invalid`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateUser_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	user := testutil.SeedUser(t, db, "Alice", "alice@example.com", "oldpass")

	body := `{"name":"Alice Updated","email":"alice.new@example.com","password":"newpass"}`
	rec := performRequest(router, http.MethodPut, fmt.Sprintf("/users/%d", user.Id), body)

	assert.Equal(t, http.StatusOK, rec.Code)

	updated := decodeUser(t, rec.Body.Bytes())
	assert.Equal(t, user.Id, updated.Id)
	assert.Equal(t, "Alice Updated", updated.Name)
	assert.Equal(t, "alice.new@example.com", updated.Email)
	assert.Equal(t, "newpass", updated.Password)
}

func TestUpdateUser_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	body := `{"name":"Ghost","email":"ghost@example.com","password":"pass"}`
	rec := performRequest(router, http.MethodPut, "/users/404", body)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `"Cannot find user with Id = 404"`, rec.Body.String())
}

func TestDeleteUser_Success(t *testing.T) {
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	user := testutil.SeedUser(t, db, "Alice", "alice@example.com", "secret")

	rec := performRequest(router, http.MethodDelete, fmt.Sprintf("/users/%d", user.Id), "")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, fmt.Sprintf(`"User with id = %d has been deleted!"`, user.Id), rec.Body.String())

	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/users/%d", user.Id), "")
	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	rec := performRequest(router, http.MethodDelete, "/users/999", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `"Cannot find user with Id = 999"`, rec.Body.String())
}

func TestUserCRUD_Flow(t *testing.T) {
	testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	createRec := performRequest(router, http.MethodPost, "/users",
		`{"name":"Flow User","email":"flow@example.com","password":"flowpass"}`)
	require.Equal(t, http.StatusOK, createRec.Code)

	created := decodeUser(t, createRec.Body.Bytes())

	listRec := performRequest(router, http.MethodGet, "/users", "")
	require.Equal(t, http.StatusOK, listRec.Code)
	users := decodeUsers(t, listRec.Body.Bytes())
	require.Len(t, users, 1)

	getRec := performRequest(router, http.MethodGet, fmt.Sprintf("/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, getRec.Code)

	updateRec := performRequest(router, http.MethodPut, fmt.Sprintf("/users/%d", created.Id),
		`{"name":"Flow Updated","email":"flow.updated@example.com","password":"newflowpass"}`)
	require.Equal(t, http.StatusOK, updateRec.Code)

	updated := decodeUser(t, updateRec.Body.Bytes())
	assert.Equal(t, "Flow Updated", updated.Name)

	deleteRec := performRequest(router, http.MethodDelete, fmt.Sprintf("/users/%d", created.Id), "")
	require.Equal(t, http.StatusOK, deleteRec.Code)

	finalListRec := performRequest(router, http.MethodGet, "/users", "")
	require.Equal(t, http.StatusOK, finalListRec.Code)
	finalUsers := decodeUsers(t, finalListRec.Body.Bytes())
	assert.Empty(t, finalUsers)
}
