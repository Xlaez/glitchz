package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRequestService interface {
	NewRequest()
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

func (g *groupRequestService) NewRequest() {}
