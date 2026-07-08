package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return val
}

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		requireEnv("DB_HOST"),
		requireEnv("DB_PORT"),
		requireEnv("DB_USER"),
		requireEnv("DB_PASSWORD"),
		requireEnv("DB_NAME"),
		requireEnv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get underlying sql.DB: %v", err))
	}

	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if maxOpen <= 0 {
		maxOpen = 25
	}
	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if maxIdle <= 0 {
		maxIdle = 10
	}
	connMaxLifeSec, _ := strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if connMaxLifeSec <= 0 {
		connMaxLifeSec = 300
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifeSec) * time.Second)

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintf("failed to auto-migrate: %v", err))
	}

	DB = db
}

func Close() {
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}
