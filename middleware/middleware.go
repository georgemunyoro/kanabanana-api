package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/munyoro/kanabanana-api/controllers"
)

func JwtAuthMiddleware(controller *controllers.GlobalController) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := controller.GetCurrentUser(c)
		if err != nil {
			println(err.Error())
			c.String(401, "Unauthorized")
			c.Abort()
			return
		}
		controller.CurrentUser = &user
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
