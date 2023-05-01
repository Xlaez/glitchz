package controllers

import (
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AuthController interface {
	RegisterUser() gin.HandlerFunc
}

type authController struct {
	s            services.AuthService
	t            services.TokenService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

type tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthController(service services.AuthService, t services.TokenService, maker token.Maker, config utils.Config, redis_client *redis.Client) AuthController {
	return &authController{
		s:            service,
		t:            t,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (a *authController) RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
