package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func ContactRoutes(router *gin.Engine, c controllers.ContactController, token_maker token.Maker) {
	auth := router.Group("/api/v1/contacts")
	auth.POST("/", c.SendReq())
}
