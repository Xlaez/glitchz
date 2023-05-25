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
	GetAllPosts() gin.HandlerFunc
	UpdatePost() gin.HandlerFunc
	DeletePost() gin.HandlerFunc
	LikePost() gin.HandlerFunc
	UnLikePost() gin.HandlerFunc
}

type postController struct {
	s            services.PostService
	n            services.NotificationService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewPostController(service services.PostService, n services.NotificationService, maker token.Maker, config utils.Config, redis_client *redis.Client) PostController {
	return &postController{
		s:            service,
		n:            n,
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

func (p *postController) GetAllPosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetAllPostsReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		filter := bson.D{{}}
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
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("resource not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"totalPosts": totalPosts, "posts": posts})
	}
}

func (p *postController) UpdatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.UpdatePostReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post id is not a valid bson objectID")))
			return
		}
		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		_, err = p.s.GetPost(bson.D{primitive.E{Key: "_id", Value: id}, {Key: "userId", Value: payload.UserID}})

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusForbidden, errorRes(errors.New("you are not the author hence don't have authority to perform this action")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		var update bson.D

		if request.Text != "" && request.Public != "" {
			var public bool
			if request.Public == "false" {
				public = false
			} else if request.Public == "true" {
				public = true
			}
			update = bson.D{{Key: "$set", Value: bson.D{{Key: "text", Value: request.Text}, {Key: "public", Value: public}, {Key: "updatedAt", Value: time.Now()}}}}
		} else if request.Text != "" {
			update = bson.D{{Key: "$set", Value: bson.D{{Key: "text", Value: request.Text}, {Key: "updatedAt", Value: time.Now()}}}}
		} else if request.Public != "" {
			var public bool
			if request.Public == "false" {
				public = false
			} else if request.Public == "true" {
				public = true
			}
			update = bson.D{{Key: "$set", Value: bson.D{{Key: "public", Value: public}, {Key: "updatedAt", Value: time.Now()}}}}
		} else {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("provide text and/or public")))
			return
		}

		if _, err = p.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("update"))
	}
}

func (p *postController) DeletePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.DeletePostReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post''s id is a valid bson objectId")))
			return
		}
		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		_, err = p.s.GetPost(bson.D{primitive.E{Key: "_id", Value: id}, {Key: "userId", Value: payload.UserID}})

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusForbidden, errorRes(errors.New("you are not the author hence don't have authority to perform this action")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err = p.s.Delete(bson.D{primitive.E{Key: "_id", Value: id}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, msgRes("deleted resource"))
	}
}

func (p *postController) LikePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.DeletePostReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post's id is not a valid objectID")))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		post, err := p.s.GetPost(filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("post not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		var update bson.D
		if post.Likes == nil {
			updateObj := []models.Likes{}
			updateObj = append(updateObj, models.Likes{
				UserID: payload.UserID,
			})
			update = bson.D{{Key: "$set", Value: bson.D{{Key: "likes", Value: updateObj}}}, {Key: "$inc", Value: bson.D{{Key: "nbLikes", Value: 1}}}}
		} else if len(post.Likes) > 0 {
			updateObj := models.Likes{
				UserID: payload.UserID,
			}
			update = bson.D{{Key: "$addToSet", Value: bson.D{{Key: "likes", Value: updateObj}}}, {Key: "$inc", Value: bson.D{{Key: "nbLikes", Value: 1}}}}
		}

		if _, err = p.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("liked!"))
	}
}

func (p *postController) UnLikePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.DeletePostReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post's id is not a valid objectID")))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		var update bson.D

		updateObj := models.Likes{
			UserID: payload.UserID,
		}

		update = bson.D{{Key: "$pull", Value: bson.D{{Key: "likes", Value: updateObj}}}, {Key: "$inc", Value: bson.D{{Key: "nbLikes", Value: -1}}}}

		if _, err = p.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("unliked!"))
	}
}

// func(p *postController)GetUserFriendsPost()gin.HandlrFunc{}
