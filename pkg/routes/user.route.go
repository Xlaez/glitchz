package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, c controllers.UserController, token_maker token.Maker) {
	auth := router.Group("/api/v1/users")
	auth.POST("/:id", c.GetUserById())
}
