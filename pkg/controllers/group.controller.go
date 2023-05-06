package controllers

import (
	"fmt"
	"glitchz/pkg/middlewares"
	"glitchz/pkg/models/group"
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
)

type GroupController interface {
	CreateGroup() gin.HandlerFunc
	GetGroupByID() gin.HandlerFunc
}

type groupController struct {
	s            services.GroupService
	request      *services.GroupRequestService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewGroupController(service services.GroupService, request *services.GroupRequestService, maker token.Maker, config utils.Config, redis_client *redis.Client) GroupController {
	return &groupController{
		s:            service,
		request:      request,
		maker:        maker,
		config:       config,
		redis_client: redis_client,
	}
}

func (g *groupController) CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.NewGroupReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		inertData := group.Group{
			ID:        primitive.NewObjectID(),
			Name:      request.Name,
			Public:    request.Public,
			CreatedAT: time.Now(),
		}
		inertData.Admins = append(inertData.Admins, group.Members{
			UserID:   payload.UserID,
			JoinedAT: time.Now(),
		})
		fmt.Print(inertData)
		for i := 0; i < len(request.Members); i++ {
			userId, err := primitive.ObjectIDFromHex(request.Members[i])
			if err != nil {
				ctx.JSON(http.StatusBadRequest, errorRes(err))
				return
			}
			inertData.Members = append(inertData.Members, group.Members{
				UserID:   userId,
				JoinedAT: time.Now(),
			})
		}

		result, err := g.s.NewGroup(inertData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (g *groupController) GetGroupByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetGroupByIDReq
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		group_id, err := primitive.ObjectIDFromHex(request.ID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}

		group, err := g.s.GetGroup(filter)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"result": group})
	}
}
