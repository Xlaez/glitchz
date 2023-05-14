package routes

import (
	"glitchz/pkg/controllers"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/services/token"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine, c controllers.PostController, token_maker token.Maker) {
	post := router.Group("/api/v1/posts").Use(middlewares.AuthMiddleWare(token_maker))
	post.POST("/", c.CreatePost())
	post.GET("/:id", c.GetPostByID())
	post.GET("/user", c.GetUserPosts())
}
