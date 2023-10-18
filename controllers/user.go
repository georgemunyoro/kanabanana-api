package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
	"github.com/munyoro/kanabanana-api/utils/token"
)

func (controller *GlobalController) CurrentUser(ctx *gin.Context) {
	user, err := controller.getCurrentUser(ctx)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An error ocurred while fetching user.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user.AsJSON(),
	})
}

func (controller *GlobalController) getCurrentUser(ctx *gin.Context) (models.User, error) {
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
