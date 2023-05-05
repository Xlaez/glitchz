package controllers

import (
	"context"
	"errors"
	"fmt"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/models"
	"glitchz/pkg/schema"
	"glitchz/pkg/services"
	"glitchz/pkg/services/password"
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
	LoginUser() gin.HandlerFunc
	RefreshAccessToken() gin.HandlerFunc
	ResetPassword() gin.HandlerFunc
	RequestPasswordReset() gin.HandlerFunc
	ConfirmPassword() gin.HandlerFunc
}

type authController struct {
	s            services.AuthService
	t            services.TokenService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
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

func (a *authController) LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.LoginReq

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := a.s.LoginUser(request.EmailOrUsername, request.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		new_user := schema.UserRes{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Pics:        user.Pics,
			Role:        user.Role,
			Bio:         user.Bio,
			LinkedIN:    user.LinkedIN,
			Github:      user.Github,
			Twitter:     user.Twitter,
			IsSuspended: user.IsSuspended,
			NBContacts:  user.NBContacts,
			NBReviews:   user.NBReviews,
		}
		tokens, err := a.t.GenerateToken(user.ID, a.config.AccessTokenDuration)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"tokens": tokens, "users": new_user})
	}
}

func (a *authController) RefreshAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.RefreshAccessTokenReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		token, err := a.maker.VerifyToken(request.Token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := a.s.GetUserById(token.UserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("user not found")))
			return
		}

		if err = a.t.DeleteToken(request.Token); err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		tokens, err := a.t.GenerateToken(user.ID, a.config.AccessTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		new_user := schema.UserRes{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Pics:        user.Pics,
			Role:        user.Role,
			Bio:         user.Bio,
			LinkedIN:    user.LinkedIN,
			Github:      user.Github,
			Twitter:     user.Twitter,
			IsSuspended: user.IsSuspended,
			NBContacts:  user.NBContacts,
			NBReviews:   user.NBReviews,
		}

		ctx.JSON(http.StatusOK, gin.H{"tokens": tokens, "user": new_user})
	}
}

func (a *authController) RequestPasswordReset() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.RequestPasswordResetReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		code, err := sendVerificationCode(ctx, a, request.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		// TODO: send code to email
		ctx.JSON(http.StatusOK, gin.H{"code": code})
	}
}

func (a *authController) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.ResetPasswordReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		email, err := a.redis_client.Get(ctx, request.Code).Result()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("6 digits code has expired or is wrong")))
			return
		}

		hashed_password, err := password.HashPassword(request.Password)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{{Key: "email", Value: email}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashed_password}}}}

		if _, err = a.s.UpdateUser(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated"))
	}
}

func (a *authController) ConfirmPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.ConfirmPassword
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		user, err := a.s.GetUserById(payload.UserID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err = password.ComparePassword(request.Password, user.Password); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		ctx.JSON(http.StatusAccepted, msgRes("corret"))
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
