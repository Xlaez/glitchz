package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func ContactRoutes(router *gin.Engine, c controllers.ContactController, token_maker token.Maker) {
	auth := router.Group("/api/v1/contacts").Use(middlewares.AuthMiddleWare(token_maker))
	auth.GET("/", c.GetContacts())
	auth.GET("/requests", c.GetContactReq())
	auth.POST("/send-request", c.SendReq())
	auth.POST("/accept-request", c.AcceptReq())
}
