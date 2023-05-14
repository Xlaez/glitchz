package controllers

import (
	"errors"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/models"
	"glitchz/pkg/schema"
	"glitchz/pkg/services"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentController interface {
	CreateComment() gin.HandlerFunc
}

type commentController struct {
	s            services.CommentService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewCommentController(service services.CommentService, maker token.Maker, config utils.Config, redis_client *redis.Client) CommentController {
	return &commentController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

// TODO: increment nbComment on post
func (c *commentController) CreateComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.NewCommentReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		post_id, err := primitive.ObjectIDFromHex(request.PostID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post's is is an invalid bsn objectId")))
			return
		}

		insertData := models.Comment{
			ID:        primitive.NewObjectID(),
			UserID:    payload.UserID,
			Content:   request.Content,
			CreatedAT: time.Now(),
			UpdatedAT: time.Now(),
			PostID:    post_id,
		}

		if request.ParentID != "" {
			parent_id, err := primitive.ObjectIDFromHex(request.ParentID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post's is is an invalid bsn objectId")))
				return
			}
			insertData.ParentID = parent_id
			if _, err := c.s.GetComment(bson.D{primitive.E{Key: "_id", Value: parent_id}}); err != nil {
				if err == mongo.ErrNoDocuments {
					ctx.JSON(http.StatusNotFound, errorRes(errors.New("comment not found")))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			filter := bson.D{primitive.E{Key: "_id", Value: parent_id}}
			update_one := bson.D{primitive.E{Key: "$inc", Value: bson.D{{Key: "nbReplies", Value: 1}}}}
			if _, err := c.s.UpdateComment(filter, update_one); err != nil {
				ctx.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
		}

		if err = c.s.NewComment(insertData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, msgRes("comment sent!"))
	}
}
