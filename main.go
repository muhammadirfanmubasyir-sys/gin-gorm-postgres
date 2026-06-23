package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	router := gin.New()
	config.Connect()
	routes.UserRoute(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
