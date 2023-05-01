package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CommentController interface {
	CreateComment() gin.HandlerFunc
}

type commentController struct {
	s            services.CommentService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewCommentController(service services.CommentService, maker token.Maker, config utils.Config, redis_client *redis.Client) CommentController {
	return &commentController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (c *commentController) CreateComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
