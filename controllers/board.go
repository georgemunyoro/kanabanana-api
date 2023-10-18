package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
)

type CreateBoardInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateBoardInput struct {
	Name           string `json:"name"`
	ListIdsInOrder string `json:"listIdsInOrder"`
}

func (controller *GlobalController) CreateBoard(ctx *gin.Context) {
	var input CreateBoardInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board := models.Board{
		Name:   input.Name,
		UserID: controller.CurrentUser.ID,
	}

	newBoard, err := board.CreateBoard(controller.Database)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create new board.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Created a new board successfully.",
		"data":    newBoard.AsJSON(true),
	})
}

func (controller *GlobalController) GetBoard(ctx *gin.Context) {
	boardId := ctx.Params.ByName("boardId")

	var board models.Board
	result := controller.Database.
		Preload("Lists").
		Preload("Lists.Cards").
		Where("id = ?", boardId).
		Where("user_id = ?", controller.CurrentUser.ID).
		First(&board)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": board.AsJSON(true),
	})
}

func (controller *GlobalController) UpdateBoard(ctx *gin.Context) {
	var input UpdateBoardInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var board models.Board
	result := controller.Database.
		Preload("Lists").
		Preload("Lists.Cards").
		Where("id = ?", ctx.Params.ByName("boardId")).
		Where("user_id = ?", controller.CurrentUser.ID).
		First(&board)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	if input.Name != board.Name && input.Name != "" {
		board.Name = input.Name
	}

	if input.ListIdsInOrder != "" &&
		input.ListIdsInOrder != board.ListIdsInOrder &&
		len(strings.Split(input.ListIdsInOrder, ",")) == len(board.Lists) {
		board.ListIdsInOrder = input.ListIdsInOrder
	}

	controller.Database.Save(board)
	ctx.JSON(http.StatusOK, gin.H{"data": board.AsJSON(true)})
}

func (controller *GlobalController) DeleteBoard(ctx *gin.Context) {
	boardId := ctx.Params.ByName("boardId")

	var board models.Board
	result := controller.Database.
		Preload("Lists").
		Preload("Lists.Cards").
		Where("id = ?", boardId).
		Where("user_id = ?", controller.CurrentUser.ID).
		First(&board)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	if controller.Database.Delete(&board).Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Board deleted successfully.",
	})
}
