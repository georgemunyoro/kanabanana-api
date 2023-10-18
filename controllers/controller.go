package controllers

import (
	"gorm.io/gorm"
)

type GlobalController struct {
	Database *gorm.DB
}
