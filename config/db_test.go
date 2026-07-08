package config

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRequireEnv_Present(t *testing.T) {
	os.Setenv("TEST_REQUIRE_ENV_KEY", "test-value")
	defer os.Unsetenv("TEST_REQUIRE_ENV_KEY")

	val := requireEnv("TEST_REQUIRE_ENV_KEY")
	assert.Equal(t, "test-value", val)
}

func TestRequireEnv_Missing_Panics(t *testing.T) {
	os.Unsetenv("TEST_REQUIRE_ENV_MISSING_KEY")

	assert.Panics(t, func() {
		requireEnv("TEST_REQUIRE_ENV_MISSING_KEY")
	})
}

func TestRequireEnv_Empty_Panics(t *testing.T) {
	os.Setenv("TEST_REQUIRE_ENV_EMPTY_KEY", "")
	defer os.Unsetenv("TEST_REQUIRE_ENV_EMPTY_KEY")

	assert.Panics(t, func() {
		requireEnv("TEST_REQUIRE_ENV_EMPTY_KEY")
	})
}

func TestConnect_MissingEnv_Panics(t *testing.T) {
	origHost := os.Getenv("DB_HOST")
	origPort := os.Getenv("DB_PORT")
	origUser := os.Getenv("DB_USER")
	origPass := os.Getenv("DB_PASSWORD")
	origName := os.Getenv("DB_NAME")
	origSSL := os.Getenv("DB_SSLMODE")
	defer func() {
		os.Setenv("DB_HOST", origHost)
		os.Setenv("DB_PORT", origPort)
		os.Setenv("DB_USER", origUser)
		os.Setenv("DB_PASSWORD", origPass)
		os.Setenv("DB_NAME", origName)
		os.Setenv("DB_SSLMODE", origSSL)
	}()

	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")

	assert.Panics(t, func() {
		Connect()
	})
}

func TestConnectWithDriver_Unsupported_Panics(t *testing.T) {
	assert.Panics(t, func() {
		ConnectWithDriver("mysql", "invalid-dsn")
	})
}

func TestConnectWithDriver_Postgres_Panics(t *testing.T) {
	assert.Panics(t, func() {
		ConnectWithDriver("postgres", "invalid-dsn")
	})
}

func TestConnectWithDialector_Success(t *testing.T) {
	origDB := DB
	defer func() { DB = origDB }()

	assert.NotPanics(t, func() {
		ConnectWithDialector(sqlite.Open(":memory:"))
	})

	require.NotNil(t, DB)
}

func TestClose_NilDB(t *testing.T) {
	orig := DB
	DB = nil
	defer func() { DB = orig }()

	assert.NotPanics(t, func() {
		Close()
	})
}

func TestClose_WithDB(t *testing.T) {
	orig := DB
	defer func() { DB = orig }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	DB = db

	assert.NotPanics(t, func() {
		Close()
	})
}

func TestConnectWithDialector_ConnPool(t *testing.T) {
	orig := DB
	origMaxOpen := os.Getenv("DB_MAX_OPEN_CONNS")
	origMaxIdle := os.Getenv("DB_MAX_IDLE_CONNS")
	origConnMaxLife := os.Getenv("DB_CONN_MAX_LIFETIME")
	defer func() {
		DB = orig
		os.Setenv("DB_MAX_OPEN_CONNS", origMaxOpen)
		os.Setenv("DB_MAX_IDLE_CONNS", origMaxIdle)
		os.Setenv("DB_CONN_MAX_LIFETIME", origConnMaxLife)
	}()

	os.Setenv("DB_MAX_OPEN_CONNS", "50")
	os.Setenv("DB_MAX_IDLE_CONNS", "20")
	os.Setenv("DB_CONN_MAX_LIFETIME", "600")

	assert.NotPanics(t, func() {
		ConnectWithDialector(sqlite.Open(":memory:"))
	})

	require.NotNil(t, DB)
	sqlDB, err := DB.DB()
	require.NoError(t, err)
	assert.Equal(t, 50, sqlDB.Stats().MaxOpenConnections)
}

func TestConnectWithDialector_DefaultPoolConfig(t *testing.T) {
	orig := DB
	origMaxOpen := os.Getenv("DB_MAX_OPEN_CONNS")
	origMaxIdle := os.Getenv("DB_MAX_IDLE_CONNS")
	origConnMaxLife := os.Getenv("DB_CONN_MAX_LIFETIME")
	defer func() {
		DB = orig
		os.Setenv("DB_MAX_OPEN_CONNS", origMaxOpen)
		os.Setenv("DB_MAX_IDLE_CONNS", origMaxIdle)
		os.Setenv("DB_CONN_MAX_LIFETIME", origConnMaxLife)
	}()

	os.Unsetenv("DB_MAX_OPEN_CONNS")
	os.Unsetenv("DB_MAX_IDLE_CONNS")
	os.Unsetenv("DB_CONN_MAX_LIFETIME")

	assert.NotPanics(t, func() {
		ConnectWithDialector(sqlite.Open(":memory:"))
	})

	require.NotNil(t, DB)
}
