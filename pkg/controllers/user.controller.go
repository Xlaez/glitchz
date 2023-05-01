package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController interface {
	GetUserById() gin.HandlerFunc
}

type userController struct {
	s            services.UserService
	maker        token.Maker
	config       utils.Config
	t_col        mongo.Collection
	redis_client *redis.Client
}

func NewUserController(service services.UserService, maker token.Maker, config utils.Config, token_col mongo.Collection, redis_client *redis.Client) UserController {
	return &userController{
		s:            service,
		maker:        maker,
		config:       config,
		t_col:        token_col,
		redis_client: redis_client,
	}
}

func (u *userController) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
