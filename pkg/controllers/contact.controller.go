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

type ContactController interface {
	SendReq() gin.HandlerFunc
	CancelReq() gin.HandlerFunc
	AcceptReq() gin.HandlerFunc
	GetContacts() gin.HandlerFunc
	BlockContact() gin.HandlerFunc
	DeleteContact() gin.HandlerFunc
	GetContactReq() gin.HandlerFunc
	UnBlockContact() gin.HandlerFunc
	GetContactByID() gin.HandlerFunc
	GetContactByUsers() gin.HandlerFunc
}

type contactController struct {
	s            services.ContactService
	n            services.NotificationService
	u            services.UserService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewContactController(service services.ContactService, n services.NotificationService, u services.UserService, maker token.Maker, config utils.Config, redis_client *redis.Client) ContactController {
	return &contactController{
		s:            service,
		n:            n,
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

		col, result, err := c.s.UpdateReq(filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

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

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		var user_id primitive.ObjectID

		if payload.UserID == result.User1 {
			user_id = result.User2
		} else if payload.UserID == result.User2 {
			user_id = result.User1
		}

		notificationChan := make(chan error)

		// go func() {
		// var err error
		contacts, err := services.GetUserContacts(ctx, col, user_id)
		if err == nil && contacts != nil {
			for i := 0; i < len(contacts); i++ {
				var contact_id primitive.ObjectID
				if contacts[i].User1 == payload.UserID {
					contact_id = contacts[i].User2
				} else if contacts[i].User2 == payload.UserID {
					contact_id = contacts[i].User1
				}
				// notificationData = append(notificationData, models.Notification{})
				_, err = c.n.NewNotification(models.Notification{
					ID:        primitive.NewObjectID(),
					UserID:    contact_id,
					Msg:       "has accepted your connection request",
					Image:     "friends",
					Seen:      false,
					CreatedAT: time.Now(),
				})
			}
			notificationChan <- err
		}

		if notificationChan != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(<-notificationChan))
			return
		}
		// }()

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

		contact_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if _, err = c.s.DeleteContact(bson.D{primitive.E{Key: "_id", Value: contact_id}, {Key: "pending", Value: true}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("request cancelled"))
	}
}

func (c *contactController) DeleteContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AcceptReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		contact_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		result, err := c.s.DeleteContact(bson.D{primitive.E{Key: "_id", Value: contact_id}, {Key: "pending", Value: false}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter_one := bson.D{primitive.E{Key: "id", Value: result.User1}}
		filter_two := bson.D{primitive.E{Key: "id", Value: result.User2}}
		update := bson.D{{Key: "$inc", Value: bson.D{{Key: "nbContacts", Value: -1}}}}

		if _, err := c.u.UpdateUser(filter_one, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if _, err := c.u.UpdateUser(filter_two, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("user is not a contact again"))
	}
}

func (c *contactController) BlockContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.BlockReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user_id, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		contact_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: contact_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "blockedIds", Value: user_id}}}}
		_, result, err := c.s.UpdateReq(filter, update)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (c *contactController) UnBlockContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AcceptReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		contact_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: contact_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "blockedIds", Value: nil}}}}
		_, result, err := c.s.UpdateReq(filter, update)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (c *contactController) GetContactByUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetContactByUser
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user_id, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		value := make([]bson.D, 0)
		value = append(value, bson.D{{Key: "user1", Value: payload.UserID}, {Key: "user2", Value: user_id}})
		value = append(value, bson.D{{Key: "user2", Value: payload.UserID}, {Key: "user1", Value: user_id}})
		filter := bson.D{{Key: "$or", Value: value}}

		contact, err := c.s.GetContact(filter, &options.FindOneOptions{})
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("contact with this user does not exist")))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": contact})
	}
}

func (c *contactController) GetContactByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetContactByID
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		contact_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		filter := bson.D{primitive.E{Key: "_id", Value: contact_id}}

		contact, err := c.s.GetContact(filter, &options.FindOneOptions{})
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("contact not found")))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": contact})
	}
}
