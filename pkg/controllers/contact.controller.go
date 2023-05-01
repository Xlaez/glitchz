package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ContactController interface {
	SendReq() gin.HandlerFunc
}

type contactController struct {
	s            services.ContactService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewContactController(service services.ContactService, maker token.Maker, config utils.Config, redis_client *redis.Client) ContactController {
	return &contactController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (c *contactController) SendReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
