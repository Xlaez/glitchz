package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type PostController interface {
	CreatePost() gin.HandlerFunc
}

type postController struct {
	s            services.PostService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewPostController(service services.PostService, maker token.Maker, config utils.Config, redis_client *redis.Client) PostController {
	return &postController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (p *postController) CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
