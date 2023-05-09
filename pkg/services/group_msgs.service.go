package services

import (
	"context"
	"glitchz/pkg/models/group"

	"go.mongodb.org/mongo-driver/mongo"
)

type GroupMsgs interface {
	SendMsg(insertBody group.GroupMsg) error
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

func (g *groupMsgs) SendMsg(insertBody group.GroupMsg) error {
	if _, err := g.col.InsertOne(g.ctx, insertBody); err != nil {
		return err
	}
	return nil
}
