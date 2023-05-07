package services

import (
	"context"
	"glitchz/pkg/models/group"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupService interface {
	NewGroup(data group.Group) (*mongo.InsertOneResult, error)
	GetGroup(filter bson.D) (group.Group, error)
	Update(filter bson.D, update bson.D) error
	GetGroups(filter bson.D, options options.FindOptions) (int64, []group.Group, error)
}

type groupService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewGroupService(col *mongo.Collection, ctx context.Context) GroupService {
	return &groupService{
		col: col,
		ctx: ctx,
	}
}

func (g *groupService) NewGroup(data group.Group) (*mongo.InsertOneResult, error) {
	result, err := g.col.InsertOne(g.ctx, data)

	if err != nil {
		return &mongo.InsertOneResult{}, err
	}
	return result, nil
}

func (g *groupService) GetGroup(filter bson.D) (group.Group, error) {
	data := group.Group{}
	if err := g.col.FindOne(g.ctx, filter, options.FindOne()).Decode(&data); err != nil {
		return group.Group{}, err
	}

	return data, nil
}

func (g *groupService) GetGroups(filter bson.D, options options.FindOptions) (int64, []group.Group, error) {
	cursor, err := g.col.Find(g.ctx, filter, &options)
	if err != nil {
		return 0, nil, err
	}

	var groups []group.Group
	if err = cursor.All(g.ctx, &groups); err != nil {
		return 0, nil, err
	}

	count, err := g.col.CountDocuments(g.ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	return count, groups, nil
}

func (g *groupService) Update(filter bson.D, update bson.D) error {
	if _, err := g.col.UpdateOne(g.ctx, filter, update); err != nil {
		return err
	}
	return nil
}
