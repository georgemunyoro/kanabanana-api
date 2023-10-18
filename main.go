package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/controllers"
	"github.com/munyoro/kanabanana-api/middleware"
	"github.com/munyoro/kanabanana-api/models"
)

func main() {
	db, err := models.ConnectDatabase()
	if err != nil {
		panic("Failed to connect to database!")
	}
	ctrl := controllers.GlobalController{Database: db}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", ctrl.Register)
			auth.POST("/login", ctrl.Login)

			me := auth.Group("/me", middleware.JwtAuthMiddleware(&ctrl))
			{
				me.GET("/", ctrl.CurrentUserHandler)
			}
		}

		board := v1.Group("/board", middleware.JwtAuthMiddleware(&ctrl))
		{
			board.POST("/", ctrl.CreateBoard)
			board.GET("/:boardId", ctrl.GetBoard)
			board.PUT("/:boardId", ctrl.UpdateBoard)
			board.DELETE("/:boardId", ctrl.DeleteBoard)

			list := board.Group("/:boardId/list")
			{
				list.POST("/", ctrl.CreateList)
				list.GET("/:listId", ctrl.GetList)
				list.PUT("/:listId", ctrl.UpdateList)
				list.DELETE("/:listId", ctrl.DeleteList)

				card := list.Group("/:listId/card")
				{
					card.POST("/", ctrl.CreateCard)
					card.GET("/:cardId", ctrl.GetCard)
					card.PUT("/:cardId", ctrl.UpdateCard)
					card.DELETE("/:cardId", ctrl.DeleteCard)
				}
			}

		}
	}

	r.Run(":8080")
}
