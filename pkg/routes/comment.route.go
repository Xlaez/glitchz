package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func CommentRoute(router *gin.Engine, c controllers.CommentController, token_maker token.Maker) {
	comment := router.Group("/api/v1/posts/comments").Use(middlewares.AuthMiddleWare(token_maker))
	comment.POST("/", c.CreateComment())
}
