package services

import (
	"context"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	GetUser(id primitive.ObjectID) models.User
}

type userService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewUserService(col *mongo.Collection, ctx context.Context) UserService {
	return &userService{
		col: col,
		ctx: ctx,
	}
}

func (a *userService) GetUser(id primitive.ObjectID) models.User {
	return models.User{}
}
