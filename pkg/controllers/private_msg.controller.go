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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PrivateMsgsController interface {
	SendMsg() gin.HandlerFunc
	DeleteMsg() gin.HandlerFunc
	SetMsgRead() gin.HandlerFunc
	AddReaction() gin.HandlerFunc
	UpdateReaction() gin.HandlerFunc
	GetRecentMsgs() gin.HandlerFunc
}

type privateMsgsController struct {
	s            services.PrivateMsgs
	c            services.ContactService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewPrivateMsgsController(service services.PrivateMsgs, c services.ContactService, maker token.Maker, config utils.Config, redis_client *redis.Client) PrivateMsgsController {
	return &privateMsgsController{
		s:            service,
		c:            c,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (p *privateMsgsController) SendMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.SendMsgReq
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		var only_text bool
		secure_url, _, err := others.UploadToCloud(ctx)
		if err != nil && request.Text == "" {
			only_text = true
		}

		if err != nil && request.Text == "" {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("provide a file or text in request body")))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		contact_id, err := primitive.ObjectIDFromHex(request.ContactID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		//NOTE: in future create a message request
		if _, err = p.c.GetContact(bson.D{primitive.E{Key: "_id", Value: contact_id}}, &options.FindOneOptions{}); err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("contact not found")))
			return
		}

		insertData := models.Msg{}
		insertData = models.Msg{
			ID:        primitive.NewObjectID(),
			ContactID: contact_id,
			Sender:    payload.UserID,
			SentAT:    time.Now(),
		}
		insertData.ReadBY = append(insertData.ReadBY, models.ReadBy{UserID: payload.UserID, ReadAT: time.Now()})
		if only_text {
			insertData.Msg = models.Message{
				Text: request.Text,
			}
		} else if request.Text == "" {
			insertData.Msg = models.Message{
				File: secure_url,
			}
		} else {
			insertData.Msg = models.Message{
				Text: request.Text,
				File: secure_url,
			}
		}

		if err = p.s.SendMsg(insertData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, msgRes("sent"))
	}
}

func (p *privateMsgsController) DeleteMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.DeleteMsgReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		msg_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		value := make([]bson.D, 0)
		value = append(value, bson.D{primitive.E{Key: "_id", Value: msg_id}})
		value = append(value, bson.D{{Key: "sender", Value: payload.UserID}})
		filter := bson.D{{Key: "$and", Value: value}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "isDeleted", Value: true}}}}

		if err = p.s.Updatemsg(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("message deleted"))
	}
}

func (p *privateMsgsController) SetMsgRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.ReadMsgReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		contact_id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		readBy := models.ReadBy{
			UserID: payload.UserID,
			ReadAT: time.Now(),
		}

		notReadBy := []models.ReadBy{}
		notReadBy = append(notReadBy, models.ReadBy{UserID: readBy.UserID})
		value := make([]bson.D, 0)
		value = append(value, bson.D{primitive.E{Key: "_id", Value: contact_id}})
		value = append(value, bson.D{{Key: "readBy", Value: bson.D{{Key: "$ne", Value: notReadBy}}}})
		filter := bson.D{{Key: "$and", Value: value}}
		update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "readBy", Value: readBy}}}}

		if err := p.s.Updatemsg(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("message updated"))
	}
}

func (p *privateMsgsController) AddReaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddReactionReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		msg_id, err := primitive.ObjectIDFromHex(request.MsgID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		ListOfReactions := make([]string, 0)
		ListOfReactions = append(ListOfReactions, "like")
		ListOfReactions = append(ListOfReactions, "dislike")
		ListOfReactions = append(ListOfReactions, "sad")
		ListOfReactions = append(ListOfReactions, "angry")
		ListOfReactions = append(ListOfReactions, "clap")

		is_reaction_valid := false

		for i := 0; i < len(ListOfReactions); i++ {
			if ListOfReactions[i] == request.Reaction {
				is_reaction_valid = true
			}
		}

		if !is_reaction_valid {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("reaction is not valid, valid reactions include: like, dislike, sad, angry and clap")))
			return
		}
		reaction := models.Reaction{
			UserID:   payload.UserID,
			Reaction: request.Reaction,
		}

		filter := bson.D{primitive.E{Key: "_id", Value: msg_id}}
		var update = bson.D{}

		Reaction := []models.Reaction{
			reaction,
		}

		msg, err := p.s.GetMsgByID(msg_id)

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("message not found")))
			return
		}

		if len(msg.Reaction) > 0 {
			update = bson.D{{Key: "$addToSet", Value: bson.D{{Key: "reaction", Value: reaction}}}}
		} else {
			update = bson.D{{Key: "$set", Value: bson.D{{Key: "reaction", Value: Reaction}}}}
		}

		if err := p.s.Updatemsg(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("reacted!"))
	}
}

// BUG HERE
func (p *privateMsgsController) UpdateReaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddReactionReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		msg_id, err := primitive.ObjectIDFromHex(request.MsgID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		ListOfReactions := make([]string, 0)
		ListOfReactions = append(ListOfReactions, "like")
		ListOfReactions = append(ListOfReactions, "dislike")
		ListOfReactions = append(ListOfReactions, "sad")
		ListOfReactions = append(ListOfReactions, "angry")
		ListOfReactions = append(ListOfReactions, "clap")

		is_reaction_valid := false

		for i := 0; i < len(ListOfReactions); i++ {
			if ListOfReactions[i] == request.Reaction {
				is_reaction_valid = true
			}
		}

		if !is_reaction_valid {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("reaction is not valid, valid reactions include: like, dislike, sad, angry and clap")))
			return
		}
		reaction := models.Reaction{
			UserID:   payload.UserID,
			Reaction: request.Reaction,
		}

		filter := bson.D{primitive.E{Key: "_id", Value: msg_id}}

		msg, err := p.s.GetMsgByID(msg_id)

		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("message not found")))
			return
		}

		for _, v := range msg.Reaction {
			if v.UserID == payload.UserID {
				v.Reaction = request.Reaction
			}
		}

		update := bson.D{{Key: "$update", Value: bson.D{{Key: "reaction", Value: reaction}}}}

		if err := p.s.Updatemsg(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("reaction updated"))
	}
}

func (p *privateMsgsController) GetRecentMsgs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
		contacts, err := p.c.GetConvsByUserId(payload.UserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		contactIds := make([]primitive.ObjectID, 0)

		for i := 0; i < len(contacts); i++ {
			contactIds = append(contactIds, contacts[i].ID)
		}

		msgs, err := p.s.GetRecentMsgs(contactIds, payload.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"results": msgs})
	}
}
