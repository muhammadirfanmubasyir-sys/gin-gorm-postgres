package routes

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/controller"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/middleware"
)

func UserRoute(router *gin.Engine) {
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.Logger())

	rateLimit := 100
	if v, err := strconv.Atoi(os.Getenv("RATE_LIMIT_PER_MINUTE")); err == nil && v > 0 {
		rateLimit = v
	}
	router.Use(middleware.RateLimit(rateLimit))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", controller.ListUser)
		v1.GET("/users/:id", controller.GetUser)
		v1.POST("/users", controller.CreateUser)
		v1.DELETE("/users/:id", controller.DeleteUser)
		v1.PUT("/users/:id", controller.UpdateUser)
	}
}
