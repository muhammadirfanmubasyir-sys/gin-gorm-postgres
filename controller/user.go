package controller

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/config"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/dto"
	"github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres/models"
)

type ValidationError struct {
	Code    string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func validate(name, email, password string) ValidationError {
	if strings.TrimSpace(name) == "" {
		return ValidationError{Code: "VALIDATION_ERROR", Message: "name is required"}
	}
	if strings.TrimSpace(email) == "" {
		return ValidationError{Code: "VALIDATION_ERROR", Message: "email is required"}
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(pattern, email); !matched {
		return ValidationError{Code: "VALIDATION_ERROR", Message: "invalid email format"}
	}

	if password == "" {
		return ValidationError{Code: "VALIDATION_ERROR", Message: "password is required"}
	}
	if len(password) < 6 {
		return ValidationError{Code: "VALIDATION_ERROR", Message: "password must be at least 6 characters"}
	}

	return ValidationError{}
}

func GetUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("user not found", "USER_NOT_FOUND"))
		return
	}

	response := dto.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

func ListUser(c *gin.Context) {
	users := []models.User{}
	config.DB.Find(&users)

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = dto.UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(responses))
}

func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("invalid request body", "INVALID_REQUEST"))
		return
	}

	if valErr := validate(req.Name, req.Email, req.Password); valErr.Code != "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(valErr.Message, valErr.Code))
		return
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to process password", "INTERNAL_ERROR"))
		return
	}

	if result := config.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to create user", "INTERNAL_ERROR"))
		return
	}

	response := dto.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

func DeleteUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("user not found", "USER_NOT_FOUND"))
		return
	}

	if result := config.DB.Delete(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to delete user", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(map[string]string{"message": "user deleted successfully"}))
}

func UpdateUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("user not found", "USER_NOT_FOUND"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("invalid request body", "INVALID_REQUEST"))
		return
	}

	if valErr := validate(req.Name, req.Email, req.Password); valErr.Code != "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(valErr.Message, valErr.Code))
		return
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to process password", "INTERNAL_ERROR"))
		return
	}

	if result := config.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to update user", "INTERNAL_ERROR"))
		return
	}

	response := dto.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}
