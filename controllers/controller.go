package controllers

import (
	"github.com/munyoro/kanabanana-api/models"
	"gorm.io/gorm"
)

type GlobalController struct {
	Database    *gorm.DB
	CurrentUser *models.User
}
