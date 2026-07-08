package routes_test

import (
	"testing"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRoute_RegistersAllEndpoints(t *testing.T) {
	router := testutil.SetupRouter()

	registered := make(map[string]bool, len(router.Routes()))
	for _, route := range router.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	expected := []string{
		"GET /api/v1/users",
		"GET /api/v1/users/:id",
		"POST /api/v1/users",
		"PUT /api/v1/users/:id",
		"DELETE /api/v1/users/:id",
	}

	for _, route := range expected {
		assert.True(t, registered[route], "expected route %s to be registered", route)
	}

	require.Len(t, registered, len(expected))
}
