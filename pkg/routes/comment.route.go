package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func CommentRoute(router *gin.Engine, c controllers.CommentController, token_maker token.Maker) {
	auth := router.Group("/api/v1/posts/comments")
	auth.POST("/", c.CreateComment())
}
