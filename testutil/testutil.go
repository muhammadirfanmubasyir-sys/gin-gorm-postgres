package testutil

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/routes"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	config.DB = db
	return db
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.UserRoute(router)
	return router
}

func SeedUser(t *testing.T, db *gorm.DB, name, email, password string) models.User {
	t.Helper()

	user := models.User{
		Name:  name,
		Email: email,
	}

	if err := user.HashPassword(password); err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	return user
}
