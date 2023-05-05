package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, c controllers.AuthController, token_maker token.Maker) {
	auth := router.Group("/api/v1/auth")
	user := router.Group("/api/v1/auth/user").Use(middlewares.AuthMiddleWare(token_maker))
	auth.POST("/register", c.RegisterUser())
	auth.POST("/verify-email", c.VerifyEmail())
	auth.POST("/login", c.LoginUser())
	auth.POST("/refresh-token", c.RefreshAccessToken())
	auth.PATCH("/reset-password", c.ResetPassword())
	auth.POST("/request-password-reset", c.RequestPasswordReset())

	user.POST("/confirm-password", c.ConfirmPassword())
}
