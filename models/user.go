package models

import (
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255;not null;" json:"name"`
	Email    string `gorm:"size:255;not null; unique;" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Boards   []Board
}

func (user *User) CreateUser(db *gorm.DB) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}

	user.Password = string(hashedPassword)
	user.Name = strings.Trim(user.Name, " ")

	err = db.Create(&user).Error
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) AsJSON() gin.H {
	boardObjectsAsJSON := make([]gin.H, len(user.Boards))
	for i, board := range user.Boards {
		boardObjectsAsJSON[i] = board.AsJSON(false)
	}

	return gin.H{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"boards": boardObjectsAsJSON,
	}
}
