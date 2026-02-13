package controller

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
)

func GetUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, "Cannot find user with Id = "+c.Param("id"))
	} else {
		c.JSON(http.StatusOK, &user)
	}

	log.Println("Get User is sucessful")
}

func ListUser(c *gin.Context) {
	users := []models.User{}

	config.DB.Find(&users)

	c.JSON(http.StatusOK, &users)

	log.Println("List User is sucessful")
}
func CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	config.DB.Create(&user)
	c.JSON(http.StatusOK, &user)

	log.Println("Create User is sucessful")
}

func DeleteUser(c *gin.Context) {
	var user models.User

	config.DB.Where("id = ?", c.Param("id")).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, "Cannot find user with Id = "+c.Param("id"))
	} else {
		config.DB.Where("id = ?", c.Param("id")).Delete(&user)
		c.JSON(http.StatusOK, "User with id = "+c.Param("id")+" has been deleted!")
	}

	log.Println("Delete User is sucessful")
}

func UpdateUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)
	if user.Id == 0 {
		c.JSON(http.StatusNotFound, "Cannot find user with Id = "+c.Param("id"))
	} else {
		c.BindJSON(&user)
		config.DB.Save(&user)
		c.JSON(http.StatusOK, &user)
	}

	log.Println("Update User is sucessful")
}
