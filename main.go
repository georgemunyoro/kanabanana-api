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

			me := auth.Group("/me").Use(middleware.JwtAuthMiddleware())
			{
				me.GET("/", ctrl.CurrentUser)
			}
		}

		board := v1.Group("/board").Use(middleware.JwtAuthMiddleware())
		{
			board.POST("/", ctrl.CreateBoard)
			board.GET("/:id", ctrl.GetBoard)
			board.PUT("/:id", ctrl.UpdateBoard)
			board.DELETE("/:id", ctrl.DeleteBoard)
		}
	}

	// product := v1.Group("/product")
	// {
	// 	product.GET("/", ctrl.GetProducts)
	// 	product.GET("/:id", ctrl.GetProduct)
	// }

	// user := v1.Group("/user")
	// user.Use(middleware.JwtAuthMiddleware())
	// {
	// 	// user.GET("/", ctrl.CurrentUser)

	// 	// cart := user.Group("/cart")
	// 	// {
	// 	// 	cart.GET("/", ctrl.GetCart)
	// 	// 	cart.POST("/:product_id", ctrl.UpsertCart)
	// 	// }
	// }

	r.Run(":8080")
}
