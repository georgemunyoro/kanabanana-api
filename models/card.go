package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	Name        string    `gorm:"not null;" json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	Attachments []Attachment
	Labels      []Label
	ListID      uint
}

func (card *Card) AsJSON() gin.H {
	attachments := make([]gin.H, len(card.Attachments))
	for i, attachment := range card.Attachments {
		attachments[i] = attachment.AsJSON()
	}

	labels := make([]gin.H, len(card.Labels))
	for i, label := range card.Labels {
		labels[i] = label.AsJSON()
	}

	return gin.H{
		"id":          card.ID,
		"name":        card.Name,
		"description": card.Description,
		"attachments": attachments,
		"labels":      labels,
		"dueDate":     card.DueDate,
		"createdAt":   card.CreatedAt,
		"updatedAt":   card.UpdatedAt,
	}
}
