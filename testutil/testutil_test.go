package testutil

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestDB_Success(t *testing.T) {
	db := SetupTestDB(t)

	assert.NotNil(t, db)
	assert.NotNil(t, config.DB)
	assert.Equal(t, db, config.DB)
}

func TestSetupTestDB_TableCreation(t *testing.T) {
	db := SetupTestDB(t)

	assert.True(t, db.Migrator().HasTable(&models.User{}))
}

func TestSetupTestDB_CanInsertData(t *testing.T) {
	db := SetupTestDB(t)

	user := models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	if err := user.HashPassword("password123"); err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	result := db.Create(&user)
	require.NoError(t, result.Error)
	assert.NotZero(t, user.Id)
}

func TestSetupTestDB_CanQueryData(t *testing.T) {
	db := SetupTestDB(t)

	user := models.User{
		Name:  "Query Test",
		Email: "query@example.com",
	}
	user.HashPassword("password123")
	db.Create(&user)

	var retrieved models.User
	db.Where("email = ?", "query@example.com").First(&retrieved)

	assert.Equal(t, user.Id, retrieved.Id)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestSetupTestDB_IsolatedInstances(t *testing.T) {
	db1 := SetupTestDB(t)

	user1 := models.User{Name: "User1", Email: "user1@example.com"}
	user1.HashPassword("pass123")
	db1.Create(&user1)

	var count1 int64
	db1.Model(&models.User{}).Count(&count1)
	require.Equal(t, int64(1), count1)

	db2 := SetupTestDB(t)

	var count2 int64
	db2.Model(&models.User{}).Count(&count2)
	require.Equal(t, int64(0), count2)
}

func TestSetupRouter_Success(t *testing.T) {
	router := SetupRouter()

	assert.NotNil(t, router)
	assert.IsType(t, &gin.Engine{}, router)
}

func TestSetupRouter_TestMode(t *testing.T) {
	_ = SetupRouter()

	require.Equal(t, "test", gin.Mode())
}

func TestSetupRouter_HasRoutes(t *testing.T) {
	router := SetupRouter()

	routes := router.Routes()
	assert.Greater(t, len(routes), 0)

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+" "+route.Path] = true
	}

	expectedRoutes := []string{
		"GET /api/v1/users",
		"GET /api/v1/users/:id",
		"POST /api/v1/users",
		"PUT /api/v1/users/:id",
		"DELETE /api/v1/users/:id",
	}

	for _, expectedRoute := range expectedRoutes {
		assert.True(t, routeMap[expectedRoute], "expected route %s not found", expectedRoute)
	}
}

func TestSeedUser_Success(t *testing.T) {
	db := SetupTestDB(t)

	user := SeedUser(t, db, "Alice", "alice@example.com", "password123")

	assert.NotZero(t, user.Id)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
}

func TestSeedUser_PasswordIsHashed(t *testing.T) {
	db := SetupTestDB(t)

	plainPassword := "mysecretpassword"
	user := SeedUser(t, db, "Bob", "bob@example.com", plainPassword)

	var retrieved models.User
	db.First(&retrieved, user.Id)

	assert.NotEqual(t, plainPassword, retrieved.Password)
	assert.NoError(t, retrieved.VerifyPassword(plainPassword))
}

func TestSeedUser_MultipleUsers(t *testing.T) {
	db := SetupTestDB(t)

	user1 := SeedUser(t, db, "User1", "user1@example.com", "pass123")
	user2 := SeedUser(t, db, "User2", "user2@example.com", "pass456")
	user3 := SeedUser(t, db, "User3", "user3@example.com", "pass789")

	var count int64
	db.Model(&models.User{}).Count(&count)
	require.Equal(t, int64(3), count)

	assert.NotEqual(t, user1.Id, user2.Id)
	assert.NotEqual(t, user2.Id, user3.Id)
}

func TestSeedUser_CanRetrieveFromDB(t *testing.T) {
	db := SetupTestDB(t)

	user := SeedUser(t, db, "Charlie", "charlie@example.com", "securepass123")

	var retrieved models.User
	db.Where("id = ?", user.Id).First(&retrieved)

	assert.Equal(t, user.Id, retrieved.Id)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestSetupTestDB_WithSeedUser(t *testing.T) {
	db := SetupTestDB(t)
	_ = SetupRouter()

	_ = SeedUser(t, db, "Alice", "alice@example.com", "pass123")
	_ = SeedUser(t, db, "Bob", "bob@example.com", "pass456")

	var users []models.User
	db.Find(&users)

	require.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "Bob", users[1].Name)
}

func TestSeedUser_DuplicateEmail(t *testing.T) {
	db := SetupTestDB(t)

	user1 := SeedUser(t, db, "Alice", "same@example.com", "pass123")
	assert.NotZero(t, user1.Id)

	_ = SeedUser(t, db, "Bob", "same@example.com", "pass456")

	var users []models.User
	db.Where("email = ?", "same@example.com").Find(&users)

	assert.Len(t, users, 2)
}

func TestSetupRouter_MultipleInstances(t *testing.T) {
	router1 := SetupRouter()
	router2 := SetupRouter()

	assert.NotNil(t, router1)
	assert.NotNil(t, router2)

	routes1 := router1.Routes()
	routes2 := router2.Routes()

	assert.Equal(t, len(routes1), len(routes2))
}
