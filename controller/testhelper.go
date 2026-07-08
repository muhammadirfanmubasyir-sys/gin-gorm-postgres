package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/repository"
)

func SetupRouterWithMock(mock repository.UserRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	ctl := NewUserController(mock)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", ctl.ListUser)
		v1.GET("/users/:id", ctl.GetUser)
		v1.POST("/users", ctl.CreateUser)
		v1.DELETE("/users/:id", ctl.DeleteUser)
		v1.PUT("/users/:id", ctl.UpdateUser)
	}

	return router
}
