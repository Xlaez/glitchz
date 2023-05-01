package controllers

import (
	"context"
	"errors"
	"fmt"
	"glitchz/pkg/models"
	"glitchz/pkg/schema"
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthController interface {
	RegisterUser() gin.HandlerFunc
	VerifyEmail() gin.HandlerFunc
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
	return func(ctx *gin.Context) {
		var request schema.AddUserReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if err := a.s.CreateUser(models.User{
			Username: request.Username,
			Password: request.Password,
			Email:    request.Email,
		}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		code, err := sendVerificationCode(ctx, a, request.Email)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		// send email here
		ctx.JSON(http.StatusCreated, gin.H{"code": code})
	}
}

func (a *authController) VerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.VerifyEmailReq

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		email, err := a.redis_client.Get(ctx, request.Code).Result()

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("code has expired or is incorrect")))
			return
		}
		filter := bson.D{{Key: "email", Value: email}}
		updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "isVerified", Value: true}}}}
		data, err := a.s.UpdateUser(filter, updateObj)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		var msg string = fmt.Sprintf(fmt.Sprint(data.ModifiedCount) + " field(s) has been modified")
		ctx.JSON(http.StatusOK, msgRes(msg))
	}
}

// In the future this function would send an email to user
func sendVerificationCode(ctx context.Context, a *authController, email string) (string, error) {
	code := utils.SixDigitsCode()
	//TODO: set expiration duration to 30 minutes
	if err := a.redis_client.Set(ctx, code, email, 0).Err(); err != nil {
		return "", err
	}

	return code, nil
}

func errorRes(err error) gin.H {
	return gin.H{"error: ": err.Error()}
}

func msgRes(msg string) gin.H {
	return gin.H{"messgae: ": msg}
}
