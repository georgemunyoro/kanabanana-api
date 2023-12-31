package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
	"github.com/munyoro/kanabanana-api/utils/token"
)

func (controller *GlobalController) CurrentUserHandler(ctx *gin.Context) {
	user := controller.CurrentUser
	if user == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An error ocurred while fetching user.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user.AsJSON(),
	})
}

func (controller *GlobalController) GetCurrentUser(ctx *gin.Context) (models.User, error) {
	currentUserId, err := token.ExtractTokenID(ctx)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	result := controller.Database.Preload("Boards").Where("id = ?", currentUserId).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}
