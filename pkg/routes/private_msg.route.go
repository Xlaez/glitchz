package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func PrivateMsgtRoutes(router *gin.Engine, c controllers.PrivateMsgsController, token_maker token.Maker) {
	auth := router.Group("/api/v1/messages").Use(middlewares.AuthMiddleWare(token_maker))
	auth.POST("/", c.SendMsg())
	auth.GET("/", c.GetRecentMsgs())
	auth.POST("/react", c.AddReaction())
	auth.PATCH("/update-reaction", c.UpdateReaction())
	auth.PATCH("/set-read/:id", c.SetMsgRead())
	auth.DELETE("/:id", c.DeleteMsg())
}
