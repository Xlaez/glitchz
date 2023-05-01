package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func PrivateMsgtRoutes(router *gin.Engine, c controllers.PrivateMsgsController, token_maker token.Maker) {
	auth := router.Group("/api/v1/messages")
	auth.POST("/", c.SendMsg())
}
