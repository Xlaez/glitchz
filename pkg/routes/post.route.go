package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine, c controllers.PostController, token_maker token.Maker) {
	auth := router.Group("/api/v1/posts")
	auth.POST("/", c.CreatePost())
}
