package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
)

type CreateBoardInput struct {
	Name string `json:"name" binding:"required"`
}

func (controller *GlobalController) CreateBoard(ctx *gin.Context) {
	user, err := controller.getCurrentUser(ctx)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An error ocurred while fetching user.",
		})
		return
	}

	var input CreateBoardInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board := models.Board{
		Name:   input.Name,
		UserID: user.ID,
	}

	newBoard, err := board.CreateBoard(controller.Database)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create new board."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Created a new board successfully.",
		"data":    newBoard.AsJSON(true),
	})
}

func (controller *GlobalController) GetBoard(ctx *gin.Context) {

}

func (controller *GlobalController) UpdateBoard(ctx *gin.Context) {

}

func (controller *GlobalController) DeleteBoard(ctx *gin.Context) {

}
