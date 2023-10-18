package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	Name           string `gorm:"not null;" json:"name"`
	CardIdsInOrder string `json:"cardIdsInOrder"`
	Cards          []Card
	BoardID        uint
}

func (list *List) CreateList(db *gorm.DB) (*List, error) {
	list.CardIdsInOrder = ""

	err := db.Create(&list).Error
	if err != nil {
		return &List{}, nil
	}
	return list, nil
}

func (list *List) AsJSON() gin.H {
	cards := make([]gin.H, len(list.Cards))
	for i, card := range list.Cards {
		cards[i] = card.AsJSON()
	}

	return gin.H{
		"id":             list.ID,
		"name":           list.Name,
		"cards":          cards,
		"cardIdsInOrder": list.CardIdsInOrder,
		"createdAt":      list.CreatedAt,
		"updatedAt":      list.UpdatedAt,
	}
}
