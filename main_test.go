package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSetupRouter(t *testing.T) {
	router := setupRouter()
	assert.NotNil(t, router)
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, "ok", body["status"])
	assert.NotNil(t, body["timestamp"])

	ts, ok := body["timestamp"].(string)
	if ok {
		parsed, err := time.Parse(time.RFC3339, ts)
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), parsed, time.Second)
	}
}

func TestHealthEndpoint_MethodNotAllowed(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSetupRouter_RoutesRegistered(t *testing.T) {
	router := setupRouter()

	routeMap := make(map[string]bool)
	for _, route := range router.Routes() {
		routeMap[route.Method+" "+route.Path] = true
	}

	assert.True(t, routeMap["GET /health"])
	assert.True(t, routeMap["GET /api/v1/users"])
	assert.True(t, routeMap["GET /api/v1/users/:id"])
	assert.True(t, routeMap["POST /api/v1/users"])
	assert.True(t, routeMap["PUT /api/v1/users/:id"])
	assert.True(t, routeMap["DELETE /api/v1/users/:id"])
}

func TestRun_DefaultAddr(t *testing.T) {
	os.Unsetenv("SERVER_ADDR")

	origDB := config.DB
	defer func() { config.DB = origDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	config.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- run()
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		t.Logf("run returned: %v", err)
	}
}

func TestRun_CustomAddr(t *testing.T) {
	os.Setenv("SERVER_ADDR", ":0")
	defer os.Unsetenv("SERVER_ADDR")

	origDB := config.DB
	defer func() { config.DB = origDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	config.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- run()
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		t.Logf("run returned: %v", err)
	}
}
