package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
	"github.com/munyoro/kanabanana-api/utils/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var REGISTRATION_ERROR = "An error ocurred while registering."
var REGISTRATION_SUCCESS = "Registration successful."
var LOGIN_ERROR = "An error ocurred while logging in."
var LOGIN_SUCCESS = "Login successful."

func (controller *GlobalController) Register(ctx *gin.Context) {
	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doesEmailExist := controller.Database.Where("email = ?", input.Email).First(&models.User{}).Error
	if doesEmailExist == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": REGISTRATION_ERROR})
		return
	}

	user := models.User{Name: input.Name, Email: input.Email, Password: input.Password}
	_, err := user.CreateUser(controller.Database)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": REGISTRATION_ERROR})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": REGISTRATION_SUCCESS,
	})
}

func (controller *GlobalController) Login(ctx *gin.Context) {
	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := controller.Database.Where("email = ?", input.Email).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": LOGIN_ERROR})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": LOGIN_ERROR})
		return
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": LOGIN_ERROR})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful.",
		"data":    token,
	})
}
