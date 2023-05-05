package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, c controllers.UserController, token_maker token.Maker) {
	auth := router.Group("/api/v1/users").Use(middlewares.AuthMiddleWare(token_maker))
	auth.GET("/", c.GetUsers())
	auth.GET("/:userId", c.GetUserById())
	auth.GET("/by-username/:username", c.GetUserByUsername())
	auth.PATCH("/update-profile", c.UpdateProfile())
	auth.PATCH("/update-pics", c.UpdatePics())
}
