package services

import (
	"context"
	"errors"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	CreateUser(data models.User) error
}

type authService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewAuthService(col *mongo.Collection, ctx context.Context) AuthService {
	return &authService{
		col: col,
		ctx: ctx,
	}
}

func (a *authService) CreateUser(data models.User) error {
	return errors.New("an error occured")
}
