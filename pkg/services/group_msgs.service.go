package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type GroupMsgs interface {
	SendMsg()
}

type groupMsgs struct {
	col *mongo.Collection
	ctx context.Context
}

func NewGroupMsgs(col *mongo.Collection, ctx context.Context) GroupMsgs {
	return &groupMsgs{
		col: col,
		ctx: ctx,
	}
}

func (g *groupMsgs) SendMsg() {}
