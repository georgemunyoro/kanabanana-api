package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	Name   string `gorm:"not null;" json:"name"`
	Color  string `json:"color"`
	CardID uint
}

func (label *Label) AsJSON() gin.H {
	return gin.H{
		"id":        label.ID,
		"name":      label.Name,
		"color":     label.Color,
		"createdAt": label.CreatedAt,
		"updatedAt": label.UpdatedAt,
	}
}
