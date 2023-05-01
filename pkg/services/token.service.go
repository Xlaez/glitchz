package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TokenService interface {
	Save()
}

type tokenService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewTokenService(col *mongo.Collection, ctx context.Context) TokenService {
	return &tokenService{
		col: col,
		ctx: ctx,
	}
}

func (t *tokenService) Save() {}
