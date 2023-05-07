package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(router *gin.Engine, c controllers.GroupController, g controllers.GroupMsgsController, token_maker token.Maker) {
	group := router.Group("/api/v1/groups").Use(middlewares.AuthMiddleWare(token_maker))
	group.POST("/", c.CreateGroup())
	group.GET("/", c.GetPublicGroups())
	group.GET("/:id", c.GetGroupByID())
	group.PUT("/add-members", c.AddUsers())
	group.PUT("/:groupId/join", c.JoinGroup())
	group.DELETE("/remove-members", c.RemoveMembers())
}
