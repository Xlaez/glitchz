package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(router *gin.Engine, c controllers.GroupController, g controllers.GroupMsgsController, token_maker token.Maker) {
	auth := router.Group("/api/v1/groups")
	auth.POST("/", c.CreateGroup())
}
