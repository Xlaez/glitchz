package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type GroupController interface {
	CreateGroup() gin.HandlerFunc
}

type groupController struct {
	s            services.GroupService
	request      *services.GroupRequestService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewGroupController(service services.GroupService, request *services.GroupRequestService, maker token.Maker, config utils.Config, redis_client *redis.Client) GroupController {
	return &groupController{
		s:            service,
		request:      request,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (p *groupController) CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
