package controllers

import (
	"errors"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/models"
	"glitchz/pkg/schema"
	"glitchz/pkg/services"
	"glitchz/pkg/services/others"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostController interface {
	CreatePost() gin.HandlerFunc
	GetPostByID() gin.HandlerFunc
	GetUserPosts() gin.HandlerFunc
}

type postController struct {
	s            services.PostService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewPostController(service services.PostService, maker token.Maker, config utils.Config, redis_client *redis.Client) PostController {
	return &postController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (p *postController) CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.NewPostReq
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		secure_url, _, err := others.UploadToCloud(ctx)

		file := ""
		if secure_url != "" && err == nil {
			file = secure_url
		}

		insertData := models.Post{
			ID:        primitive.NewObjectID(),
			UserID:    payload.UserID,
			Text:      request.Text,
			FIle:      file,
			Public:    request.Public,
			CreatedAT: time.Now(),
			UpdatedAT: time.Now(),
		}

		if err = p.s.NewPost(insertData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ctx.JSON(http.StatusCreated, msgRes("post created"))
	}
}

func (p *postController) GetPostByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetPostByIDReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post's id is not a valid bson object Id")))
			return
		}

		post, err := p.s.GetPost(bson.D{primitive.E{Key: "_id", Value: id}})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("post not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"result": post})
	}
}

func (p *postController) GetUserPosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetPostsByUserIDReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user_id, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("user's id is not a valid bson objectID")))
			return
		}

		filter := bson.D{{Key: "userId", Value: user_id}}
		counter := int64(1)
		page := int64(request.Page)
		limit := int64(request.Limit)
		skip := (page - counter) * limit

		options := &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		}

		totalPosts, posts, err := p.s.GetPosts(filter, options)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("this user does not have any posts")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"totalPosts": totalPosts, "posts": posts})
	}
}
