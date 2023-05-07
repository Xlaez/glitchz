package controllers

import (
	"errors"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupController interface {
	SendRequestToPrivateGroup() gin.HandlerFunc
	GetGroupRequests() gin.HandlerFunc
	RemoveMembers() gin.HandlerFunc
	GetPublicGroups() gin.HandlerFunc
	RemoveAdmins() gin.HandlerFunc
	GetGroupByID() gin.HandlerFunc
	CreateGroup() gin.HandlerFunc
	UnBlockUser() gin.HandlerFunc
	BlockUser() gin.HandlerFunc
	JoinGroup() gin.HandlerFunc
	AddUsers() gin.HandlerFunc
	AddAdmin() gin.HandlerFunc
}

type groupController struct {
	s            services.GroupService
	request      services.GroupRequestService
	maker        token.Maker
	config       utils.Config
	redis_client *redis.Client
}

func NewGroupController(service services.GroupService, request services.GroupRequestService, maker token.Maker, config utils.Config, redis_client *redis.Client) GroupController {
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

func (g *groupController) GetPublicGroups() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetGroupsReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		filter := bson.D{{Key: "public", Value: true}}
		counter := int64(1)
		skip := (request.Page - counter) * request.Limit
		options := &options.FindOptions{
			Limit: &request.Limit,
			Skip:  &skip,
		}

		count, groups, err := g.s.GetGroups(filter, *options)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"count": count, "groups": groups})
	}
}

func (g *groupController) AddUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddMembersToGroup
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		members := []group.Members{}
		for i := 0; i < len(request.IDs); i++ {
			members = append(members, group.Members{
				UserID:   request.IDs[i],
				JoinedAT: time.Now(),
			})
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "members", Value: bson.D{{Key: "$each", Value: members}}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("added"))
	}
}

func (g *groupController) JoinGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.JoinGroup
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		_, err = g.s.GetGroup(bson.D{primitive.E{Key: "_id", Value: group_id}, {Key: "public", Value: true}})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotAcceptable, errorRes(errors.New("you cannot join a private group using invite link")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "members", Value: group.Members{UserID: payload.UserID, JoinedAT: time.Now()}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("added"))
	}
}

// TODO: fix bug that refususes to delete users
func (g *groupController) RemoveMembers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddMembersToGroup
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		members := []group.Members{}
		for i := 0; i < len(request.IDs); i++ {
			members = append(members, group.Members{
				UserID: request.IDs[i],
			})
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$pull", Value: bson.D{{Key: "members", Value: bson.D{{Key: "$in", Value: members}}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("removed"))
	}
}

func (g *groupController) BlockUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddMembersToGroup
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		members := []primitive.ObjectID{}
		for i := 0; i < len(request.IDs); i++ {
			members = append(members, request.IDs[i])
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "blockedIds", Value: bson.D{{Key: "$each", Value: members}}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("blocked!"))
	}
}

func (g *groupController) UnBlockUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddMembersToGroup
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		members := []primitive.ObjectID{}
		for i := 0; i < len(request.IDs); i++ {
			members = append(members, request.IDs[i])
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$pull", Value: bson.D{{Key: "blockedIds", Value: bson.D{{Key: "$in", Value: members}}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("unblocked!"))
	}
}

func (g *groupController) AddAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddMembersToGroup
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		members := []group.Members{}
		for i := 0; i < len(request.IDs); i++ {
			members = append(members, group.Members{
				UserID:   request.IDs[i],
				JoinedAT: time.Now(),
			})
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "admins", Value: bson.D{{Key: "$each", Value: members}}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("added"))
	}
}

func (g *groupController) RemoveAdmins() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AddMembersToGroup
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		members := []group.Members{}
		for i := 0; i < len(request.IDs); i++ {
			members = append(members, group.Members{
				UserID: request.IDs[i],
			})
		}

		filter := bson.D{primitive.E{Key: "_id", Value: group_id}}
		update := bson.D{{Key: "$pull", Value: bson.D{{Key: "admins", Value: bson.D{{Key: "$in", Value: members}}}}}}

		if err = g.s.Update(filter, update); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, msgRes("removed"))
	}
}

func (g *groupController) SendRequestToPrivateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.SendRequestReq
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		if err = g.request.NewRequest(payload.UserID, group_id, request.Msg); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err = g.s.Update(bson.D{primitive.E{Key: "_id", Value: group_id}}, bson.D{{Key: "$inc", Value: bson.D{{Key: "nbRequests", Value: 1}}}}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, msgRes("message sent"))
	}
}

func (g *groupController) GetGroupRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.GetGroupRequestsReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		counter := int64(1)
		skip := (request.Page - counter) * request.Limit
		options := &options.FindOptions{
			Limit: &request.Limit,
			Skip:  &skip,
		}

		group_id, err := primitive.ObjectIDFromHex(request.GroupID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "groupId", Value: group_id}}

		requests, count, err := g.request.GetRequests(filter, *options)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"count": count, "requests": requests})
	}
}

func (g *groupController) GetRequestByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request schema.AcceptReq
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		// request, err := g.request.GetRequestByID(request.ID)
	}
}
