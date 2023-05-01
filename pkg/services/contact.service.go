package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ContactService interface {
	SendReq()
}

type contactService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewContactService(col *mongo.Collection, ctx context.Context) ContactService {
	return &contactService{
		col: col,
		ctx: ctx,
	}
}

func (c *contactService) SendReq() {}
