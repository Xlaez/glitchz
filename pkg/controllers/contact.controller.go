package controllers

import (
	"fmt"
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

type ContactController interface {
	SendReq() gin.HandlerFunc
	AcceptReq() gin.HandlerFunc
	GetContacts() gin.HandlerFunc
	GetContactReq() gin.HandlerFunc
}

type contactController struct {
	s            services.ContactService
	u            services.UserService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewContactController(service services.ContactService, u services.UserService, maker token.Maker, config utils.Config, redis_client *redis.Client) ContactController {
	return &contactController{
		s:            service,
		u:            u,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (c *contactController) SendReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.SendReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id := primitive.NewObjectID()
		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		user_id, err := primitive.ObjectIDFromHex(request.UserID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		insertObj := models.Contact{
			ID:      id,
			User1:   user_id,
			User2:   payload.UserID,
			Pending: true,
			SentAT:  time.Now(),
		}

		//TODO: use the id returned from this to send notification
		if _, err = c.s.SendReq(insertObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("sent"))
	}
}

func (c *contactController) AcceptReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AcceptReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		request_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: request_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "pending", Value: false}, {Key: "acceptedAt", Value: time.Now()}}}}

		result, err := c.s.UpdateReq(filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		fmt.Print(result)
		userFilter1 := bson.D{primitive.E{Key: "id", Value: result.User1}}
		userFilter2 := bson.D{primitive.E{Key: "id", Value: result.User2}}
		userUpdate := bson.D{{Key: "$inc", Value: bson.D{{Key: "nbContacts", Value: 1}}}}

		if _, err := c.u.UpdateUser(userFilter1, userUpdate); err != nil {
			ctx.JSONP(http.StatusInternalServerError, errorRes(err))
			return
		}

		if _, err := c.u.UpdateUser(userFilter2, userUpdate); err != nil {
			ctx.JSONP(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("you are now contacts"))
	}
}

func (c *contactController) GetContacts() gin.HandlerFunc {
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

		options := &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		value := make([]primitive.D, 0)
		value = append(value, bson.D{{Key: "user1", Value: payload.UserID}, {Key: "pending", Value: false}})
		value = append(value, bson.D{{Key: "user2", Value: payload.UserID}, {Key: "pending", Value: false}})

		contacts, total_contacts, err := c.s.GetContacts(bson.D{{Key: "$or", Value: value}}, options)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"results": contacts, "totalContacts": total_contacts})
	}
}

func (c *contactController) GetContactReq() gin.HandlerFunc {
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

		options := &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		value := make([]primitive.D, 0)
		value = append(value, bson.D{{Key: "user1", Value: payload.UserID}, {Key: "pending", Value: true}})

		contacts, total_contacts, err := c.s.GetContacts(bson.D{{Key: "$or", Value: value}}, options)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"results": contacts, "totalRequests": total_contacts})
	}
}

func (c *contactController) CancelReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AcceptReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

	}
}
