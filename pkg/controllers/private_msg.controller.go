package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type PrivateMsgsController interface {
	SendMsg() gin.HandlerFunc
}

type privateMsgsController struct {
	s            services.PrivateMsgs
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewPrivateMsgsController(service services.PrivateMsgs, maker token.Maker, config utils.Config, redis_client *redis.Client) PrivateMsgsController {
	return &privateMsgsController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (p *privateMsgsController) SendMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
