package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController interface {
	RegisterUser() gin.HandlerFunc
}

type authController struct {
	s            services.AuthService
	maker        token.Maker
	config       utils.Config
	t_col        mongo.Collection
	redis_client *redis.Client
}

type tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthController(service services.AuthService, maker token.Maker, config utils.Config, token_col mongo.Collection, redis_client *redis.Client) AuthController {
	return &authController{
		s:            service,
		maker:        maker,
		config:       config,
		t_col:        token_col,
		redis_client: redis_client,
	}
}

func (a *authController) RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
