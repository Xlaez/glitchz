package services

import (
	"context"
	"glitchz/pkg/models/group"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupRequestService interface {
	NewRequest(userId primitive.ObjectID, groupId primitive.ObjectID, msg string) error
	GetRequests(filter bson.D, option options.FindOptions) ([]group.Request, int64, error)
	GetRequestByID(id primitive.ObjectID) (group.Request, error)
	DeleteRequest(filter bson.D) error
}

type groupRequestService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewGroupRequestService(col *mongo.Collection, ctx context.Context) GroupRequestService {
	return &groupRequestService{
		col: col,
		ctx: ctx,
	}
}

func (g *groupRequestService) NewRequest(userId primitive.ObjectID, groupId primitive.ObjectID, msg string) error {
	request := group.Request{
		ID:      primitive.NewObjectID(),
		GroupID: groupId,
		UserID:  userId,
		SentAT:  time.Now(),
		Msg:     msg,
	}

	if _, err := g.col.InsertOne(g.ctx, request, options.InsertOne()); err != nil {
		return err
	}

	return nil
}

func (g *groupRequestService) GetRequests(filter bson.D, option options.FindOptions) ([]group.Request, int64, error) {
	requests := []group.Request{}

	cursor, err := g.col.Find(g.ctx, filter, &option)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(g.ctx, &requests); err != nil {
		return nil, 0, err
	}

	count, err := g.col.CountDocuments(g.ctx, filter, options.Count())
	if err != nil {
		return nil, 0, err
	}

	return requests, count, nil
}

func (g *groupRequestService) GetRequestByID(request_id primitive.ObjectID) (group.Request, error) {
	request := group.Request{}
	if err := g.col.FindOne(g.ctx, bson.D{primitive.E{Key: "_id", Value: request_id}}, options.FindOne()).Decode(&request); err != nil {
		return group.Request{}, err
	}

	return request, nil
}

func (g *groupRequestService) DeleteRequest(filter bson.D) error {
	if _, err := g.col.DeleteOne(g.ctx, filter); err != nil {
		return err
	}
	return nil
}
