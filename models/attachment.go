package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	FileName string `gorm:"not null;" json:"name"`
	Url      string `json:"url"`
	MimeType string `json:"mimeType"`
	CardID   uint
}

func (attachment *Attachment) AsJSON() gin.H {
	return gin.H{
		"id":        attachment.ID,
		"filename":  attachment.FileName,
		"url":       attachment.FileName,
		"mimeType":  attachment.MimeType,
		"createdAt": attachment.CreatedAt,
	}
}
