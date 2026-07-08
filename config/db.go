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

	ConnectWithDialector(postgres.Open(dsn))
}

func ConnectWithDriver(driver, dsn string) {
	var dialector gorm.Dialector
	switch driver {
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		panic(fmt.Sprintf("unsupported database driver: %s", driver))
	}

	ConnectWithDialector(dialector)
}

func ConnectWithDialector(dialector gorm.Dialector) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	configureConnectionPool(db)

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintf("failed to auto-migrate: %v", err))
	}

	DB = db
}

func configureConnectionPool(db *gorm.DB) {
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
}

func Close() {
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}
