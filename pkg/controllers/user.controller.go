package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type UserController interface {
	GetUserById() gin.HandlerFunc
}

type userController struct {
	s            services.UserService
	t            services.TokenService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewUserController(service services.UserService, t services.TokenService, maker token.Maker, config utils.Config, redis_client *redis.Client) UserController {
	return &userController{
		s:            service,
		t:            t,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (u *userController) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
