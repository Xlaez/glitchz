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
	auth.GET("/by-id", c.GetContactByID())
	auth.GET("/by-users", c.GetContactByUsers())
	auth.GET("/requests", c.GetContactReq())
	auth.POST("/block", c.BlockContact())
	auth.POST("/unblock", c.UnBlockContact())
	auth.POST("/send-request", c.SendReq())
	auth.POST("/cancel-request", c.CancelReq())
	auth.POST("/remove", c.DeleteContact())
	auth.POST("/accept-request", c.AcceptReq())
}
