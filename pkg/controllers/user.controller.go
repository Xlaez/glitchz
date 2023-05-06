package controllers

import (
	"errors"
	"fmt"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/schema"
	"glitchz/pkg/services"
	"glitchz/pkg/services/others"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserController interface {
	GetUserById() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	GetUserByUsername() gin.HandlerFunc
	UpdateProfile() gin.HandlerFunc
	UpdatePics() gin.HandlerFunc
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
	return func(ctx *gin.Context) {
		var request schema.GetUserByIdReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		userId, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		user, err := u.s.GetUser(bson.D{primitive.E{Key: "id", Value: userId}})

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
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

		ctx.JSON(http.StatusOK, gin.H{"user": new_user})
	}
}

func (u *userController) GetUserByUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetUserByUsername
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		username := fmt.Sprint("@" + request.Username)
		user, err := u.s.GetUser(bson.D{{Key: "username", Value: username}})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("user not found")))
				return
			}
			ctx.JSON(http.StatusNotFound, errorRes(err))
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

		ctx.JSON(http.StatusOK, gin.H{"user": new_user})
	}
}

func (u *userController) UpdateProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.UpdateProfileReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload, _ := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		fmt.Print(payload)

		filter := bson.D{primitive.E{Key: "id", Value: payload.UserID}}

		var updateObj bson.D
		if request.Bio != "" && request.Github != "" && request.LinkedIN != "" && request.Twitter != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "bio", Value: request.Bio}, {Key: "github", Value: request.Github}, {Key: "linkedIn", Value: request.LinkedIN}, {Key: "twitter", Value: request.Twitter}}}}
		} else if request.Github != "" && request.LinkedIN != "" && request.Twitter != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "github", Value: request.Github}, {Key: "linkedIn", Value: request.LinkedIN}, {Key: "twitter", Value: request.Twitter}}}}
		} else if request.LinkedIN != "" && request.Twitter != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "linkedIn", Value: request.LinkedIN}, {Key: "twitter", Value: request.Twitter}}}}
		} else if request.Github != "" && request.Twitter != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "github", Value: request.Github}, {Key: "twitter", Value: request.Twitter}}}}
		} else if request.LinkedIN != "" && request.Github != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "github", Value: request.Github}, {Key: "linkedIn", Value: request.LinkedIN}}}}
		} else if request.Bio != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "bio", Value: request.Bio}}}}
		} else if request.Github != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "github", Value: request.Github}}}}
		} else if request.LinkedIN != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "linkedIn", Value: request.LinkedIN}}}}
		} else if request.Twitter != "" {
			updateObj = bson.D{{Key: "$set", Value: bson.D{{Key: "twitter", Value: request.Twitter}}}}
		}

		result, err := u.s.UpdateUser(filter, updateObj)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result})
	}
}

// type UpdatePics struct{}

func (u *userController) UpdatePics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var request UpdatePics
		// if err := ctx.ShouldBind(&request); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, errorRes(err))
		// 	return
		// }
		secure_url, _, err := others.UploadToCloud(ctx)
		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		filter := bson.D{primitive.E{Key: "id", Value: payload.UserID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "pics", Value: secure_url}}}}
		result, err := u.s.UpdateUser(filter, update)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (u *userController) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetUsersReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		counter := int64(1)
		page := int64(request.Page)
		limit := int64(request.Limit)
		skip := (page - counter) * limit
		filter := bson.D{}
		options := &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		}

		result, totalDocs, err := u.s.GetUsers(filter, options)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusInternalServerError, errorRes(errors.New("resources not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result, "totalUsers": totalDocs})
	}
}
