package routes

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/controller"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/middleware"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/repository"
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

	repo := repository.NewUserRepositoryGorm(config.DB)
	ctl := controller.NewUserController(repo)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", ctl.ListUser)
		v1.GET("/users/:id", ctl.GetUser)
		v1.POST("/users", ctl.CreateUser)
		v1.DELETE("/users/:id", ctl.DeleteUser)
		v1.PUT("/users/:id", ctl.UpdateUser)
	}
}
