package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/models"
)

type CreateCardInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateCardInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (controller *GlobalController) CreateCard(ctx *gin.Context) {
	var list models.List
	result := controller.Database.
		Preload("Cards").
		Where("id = ?", ctx.Params.ByName("listId")).
		Where("board_id = ?", ctx.Params.ByName("boardId")).
		Find(&list)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Unable to find list."})
		return
	}

	var input CreateCardInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCard := models.Card{
		Name:        input.Name,
		Description: input.Description,
		ListID:      list.ID,
	}

	err := controller.Database.
		Model(&list).
		Association("Cards").
		Append(&newCard)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create new card.",
		})
		return
	}

	result = controller.Database.Save(list)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create new card.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Created new card successfully.",
		"data":    newCard.AsJSON(),
	})
}

func (controller *GlobalController) GetCard(ctx *gin.Context) {
	var card models.Card
	result := controller.Database.
		Preload("Attachments").
		Preload("Labels").
		Where("id = ?", ctx.Params.ByName("cardId")).
		Where("list_id = ?", ctx.Params.ByName("listId")).
		First(&card)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": card.AsJSON(),
	})
}

func (controller *GlobalController) UpdateCard(ctx *gin.Context) {
	var input UpdateCardInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var card models.Card
	result := controller.Database.
		Preload("Attachments").
		Preload("Labels").
		Where("id = ?", ctx.Params.ByName("cardId")).
		Where("list_id", ctx.Params.ByName("listId")).
		First(&card)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	if input.Name != card.Name && input.Name != "" {
		card.Name = input.Name
	}

	if input.Description != card.Description {
		card.Description = input.Description
	}

	controller.Database.Save(card)
	ctx.JSON(http.StatusOK, gin.H{"data": card.AsJSON()})
}

func (controller *GlobalController) DeleteCard(ctx *gin.Context) {
	var card models.Card
	result := controller.Database.
		Preload("Attachments").
		Preload("Labels").
		Where("id = ?", ctx.Params.ByName("cardId")).
		Where("list_id", ctx.Params.ByName("listId")).
		First(&card)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	if controller.Database.Delete(&card).Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "An unexpected error ocurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Card deleted successfully.",
	})
}
