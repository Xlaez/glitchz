package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, c controllers.AuthController, token_maker token.Maker) {
	auth := router.Group("/api/v1/auth")
	auth.POST("/register", c.RegisterUser())
}
