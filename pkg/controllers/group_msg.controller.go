package controllers

import (
	"errors"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/models"
	"glitchz/pkg/models/group"
	"glitchz/pkg/schema"
	"glitchz/pkg/services"
	"glitchz/pkg/services/others"
	"glitchz/pkg/services/token"
	"glitchz/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupMsgsController interface {
	SendMsg() gin.HandlerFunc
}

type groupMsgsController struct {
	s            services.GroupMsgs
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewGroupMsgsController(service services.GroupMsgs, maker token.Maker, config utils.Config, redis_client *redis.Client) GroupMsgsController {
	return &groupMsgsController{
		s:            service,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (g *groupMsgsController) SendMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.NewGroupMsgReq
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		secure_url, _, err := others.UploadToCloud(ctx)

		var message models.Message

		if secure_url != "" && err == nil && request.Text != "" {
			message = models.Message{
				Text: request.Text,
				File: secure_url,
			}
		} else if secure_url != "" {
			message = models.Message{
				File: secure_url,
			}
		} else if request.Text != "" {
			message = models.Message{
				Text: request.Text,
			}
		} else {
			ctx.JSON(http.StatusBadRequest, errorRes(errors.New("attach file or/and text to request body")))
			return
		}

		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		insertData := group.GroupMsg{
			ID:      primitive.NewObjectID(),
			GroupID: group_id,
			Sender:  payload.UserID,
			Msg:     message,
		}

		insertData.ReadBY = append(insertData.ReadBY, models.ReadBy{
			UserID: payload.UserID,
			ReadAT: time.Now(),
		})
		if err = g.s.SendMsg(insertData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("message sent"))
	}
}
