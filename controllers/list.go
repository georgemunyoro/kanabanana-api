package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
)

type CreateListInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateListInput struct {
	Name           string `json:"name"`
	CardIdsInOrder string `json:"cardIdsInOrder"`
}

func (controller *GlobalController) CreateList(ctx *gin.Context) {
	var board models.Board
	result := controller.Database.
		Preload("Lists").
		Where("id = ?", ctx.Params.ByName("boardId")).
		Where("user_id = ?", controller.CurrentUser.ID).
		Find(&board)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Unable to find board."})
		return
	}

	var input CreateListInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newList := models.List{
		Name:    input.Name,
		BoardID: board.ID,
	}

	err := controller.Database.
		Model(&board).
		Association("Lists").
		Append(&newList)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create new list.",
		})
		return
	}

	result = controller.Database.Save(board)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create new list.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Created new list.",
		"data":    newList.AsJSON(),
	})
}

func (controller *GlobalController) GetList(ctx *gin.Context) {
	var list models.List
	result := controller.Database.
		Preload("Cards").
		Where("id = ?", ctx.Params.ByName("listId")).
		Where("board_id = ?", ctx.Params.ByName("boardId")).
		First(&list)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": list.AsJSON(),
	})
}

func (controller *GlobalController) UpdateList(ctx *gin.Context) {
	var input UpdateListInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var list models.List
	result := controller.Database.
		Preload("Cards").
		Where("id = ?", ctx.Params.ByName("listId")).
		Where("board_id", ctx.Params.ByName("boardId")).
		First(&list)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	if input.Name != list.Name && input.Name != "" {
		list.Name = input.Name
	}

	if input.CardIdsInOrder != "" &&
		input.CardIdsInOrder != list.CardIdsInOrder &&
		len(strings.Split(input.CardIdsInOrder, ",")) == len(list.Cards) {
		list.CardIdsInOrder = input.CardIdsInOrder
	}

	controller.Database.Save(list)
	ctx.JSON(http.StatusOK, gin.H{"data": list.AsJSON()})
}

func (controller *GlobalController) DeleteList(ctx *gin.Context) {
	var list models.List
	result := controller.Database.
		Preload("Cards").
		Where("id = ?", ctx.Params.ByName("listId")).
		Where("board_id = ?", ctx.Params.ByName("boardId")).
		First(&list)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
	}

	if controller.Database.Delete(&list).Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "List deleted successfully.",
	})
}
