package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(router *gin.Engine, c controllers.GroupController, g controllers.GroupMsgsController, token_maker token.Maker) {
	group := router.Group("/api/v1/groups").Use(middlewares.AuthMiddleWare(token_maker))
	group.GET("/", c.GetPublicGroups())
	group.GET("/requests", c.GetGroupRequests())
	group.GET("/requests/:id", c.GetRequestByID())
	group.GET("/:id", c.GetGroupByID())
	group.POST("/", c.CreateGroup())
	group.POST("/add-members", c.AddUsers())
	group.POST("/add-admins", c.AddAdmin())
	group.POST("/:groupId/join", c.JoinGroup())
	group.POST("/send-request", c.SendRequestToPrivateGroup())
	group.POST("/accept-request/:id", c.AcceptRequest())
	group.POST("/cancel-request/:id", c.CancelRequest())
	group.POST("/block", c.BlockUser())
	group.DELETE("/:id", c.DeleteGroup())
	group.DELETE("/unblock", c.UnBlockUser())
	group.DELETE("/remove-members", c.RemoveMembers())
	group.DELETE("/remove-admins", c.RemoveAdmins())
}
