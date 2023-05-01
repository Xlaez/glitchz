package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type GroupMsgsController interface {
	SendMsg() gin.HandlerFunc
}

type groupMsgsController struct {
	s            services.GroupMsgs
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewGroupMsgsController(service services.GroupMsgs, maker token.Maker, config utils.Config, redis_client *redis.Client) GroupMsgsController {
	return &groupMsgsController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (g *groupMsgsController) SendMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
