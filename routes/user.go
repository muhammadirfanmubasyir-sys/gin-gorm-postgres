package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/controller"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", controller.ListUser)
	router.GET("/:id", controller.GetUser)
	router.POST("/", controller.CreateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.PUT("/:id", controller.UpdateUser)

}
