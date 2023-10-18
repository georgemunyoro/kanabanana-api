package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Board struct {
	gorm.Model
	Name           string `gorm:"not null;" json:"name"`
	ListIdsInOrder string `json:"listIdsInOrder"`
	Lists          []List
	UserID         uint
}

func (board *Board) CreateBoard(db *gorm.DB) (*Board, error) {
	board.ListIdsInOrder = ""

	err := db.Create(&board).Error
	if err != nil {
		return &Board{}, nil
	}
	return board, nil
}

func (board *Board) AsJSON(showLists bool) gin.H {
	listObjectsAsJSON := make([]gin.H, len(board.Lists))
	for i, list := range board.Lists {
		listObjectsAsJSON[i] = list.AsJSON()
	}

	if showLists {
		return gin.H{
			"id":             board.ID,
			"name":           board.Name,
			"lists":          listObjectsAsJSON,
			"listIdsInOrder": board.ListIdsInOrder,
			"createdAt":      board.CreatedAt,
			"updatedAt":      board.UpdatedAt,
		}
	} else {
		return gin.H{
			"id":        board.ID,
			"name":      board.Name,
			"createdAt": board.CreatedAt,
			"updatedAt": board.UpdatedAt,
		}
	}
}
