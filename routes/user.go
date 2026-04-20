package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/controller"
)

func UserRoute(router *gin.Engine) {
	router.GET("/users", controller.ListUser)
	router.GET("/users/:id", controller.GetUser)
	router.POST("/users", controller.CreateUser)
	router.DELETE("/users/:id", controller.DeleteUser)
	router.PUT("/users/:id", controller.UpdateUser)

}
