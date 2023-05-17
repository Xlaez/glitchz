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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentController interface {
	GetComments() gin.HandlerFunc
	CreateComment() gin.HandlerFunc
	UpdateComment() gin.HandlerFunc
	DeleteComment() gin.HandlerFunc
	GetCommentReplies() gin.HandlerFunc
}

type commentController struct {
	s            services.CommentService
	p            services.PostService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewCommentController(service services.CommentService, p services.PostService, maker token.Maker, config utils.Config, redis_client *redis.Client) CommentController {
	return &commentController{
		s:            service,
		p:            p,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

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

		if _, err = c.p.Update(bson.D{primitive.E{Key: "_id", Value: post_id}}, bson.D{{Key: "$inc", Value: bson.D{{Key: "nbComments", Value: 1}}}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, msgRes("comment sent!"))
	}
}

func (c *commentController) GetComments() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetCommentsReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		post_id, err := primitive.ObjectIDFromHex(request.PostID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("post's id is not a valid objectID")))
			return
		}

		counter := int64(1)
		page := int64(request.Page)
		limit := int64(request.Limit)
		skip := (page - counter) * limit

		options := &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		}

		filter := bson.D{primitive.E{Key: "postId", Value: post_id}}
		comments, err := c.s.GetComments(filter, options)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("no comments for this post")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": comments})
	}
}

func (c *commentController) GetCommentReplies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetCommentsRepliesReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		comment_id, err := primitive.ObjectIDFromHex(request.CommentID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("comment's id is not a valid objectID")))
			return
		}

		counter := int64(1)
		page := int64(request.Page)
		limit := int64(request.Limit)
		skip := (page - counter) * limit

		options := &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		}

		filter := bson.D{primitive.E{Key: "parentId", Value: comment_id}}
		comments, err := c.s.GetComments(filter, options)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("no replies for this comment")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": comments})
	}
}

func (c *commentController) UpdateComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.UpdateCommentReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		comment_id, err := primitive.ObjectIDFromHex(request.CommentID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("comment's id is an invalid objectID")))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: comment_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: request.Text},
			{Key: "updatedAt", Value: time.Now()}}}}

		if _, err = c.s.UpdateComment(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("updated!"))
	}
}

func (c *commentController) DeleteComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.DeleteCommentReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		comment_id, err := primitive.ObjectIDFromHex(request.ID)

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("comment's id is not a valid objectID")))
			return
		}

		comment, err := c.s.GetComment(bson.D{primitive.E{Key: "_id", Value: comment_id}})

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("comment's id does not match any comment")))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		if payload.UserID != comment.UserID {
			ctx.JSON(http.StatusUnauthorized, errorRes(errors.New("user is not authorized to delete this resource")))
			return
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("comment's id is an invalid objectID")))
			return
		}

		if primitive.IsValidObjectID(comment.ParentID.Hex()) {
			if _, err = c.s.UpdateComment(bson.D{primitive.E{Key: "_id", Value: comment.ParentID}}, bson.D{{Key: "$inc", Value: bson.D{{Key: "nbReplies", Value: -1}}}}); err != nil {
				ctx.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
		}

		if err = c.s.DeleteOneComment(bson.D{primitive.E{Key: "_id", Value: comment_id}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err = c.s.DeleteManyComment(bson.D{primitive.E{Key: "parentId", Value: comment_id}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if _, err = c.p.Update(bson.D{primitive.E{Key: "_id", Value: comment.PostID}}, bson.D{{Key: "$inc", Value: bson.D{{Key: "nbComments", Value: -1}}}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("deleted!"))
	}
}
