package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PrivateMsgs interface {
	SendMsg()
}

type privateMsgs struct {
	col *mongo.Collection
	ctx context.Context
}

func NewPrivateMsgs(col *mongo.Collection, ctx context.Context) PrivateMsgs {
	return &privateMsgs{
		col: col,
		ctx: ctx,
	}
}

func (p *privateMsgs) SendMsg() {}
