package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type GroupService interface {
	NewGroup()
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

func (g *groupService) NewGroup() {}
