package services

import (
	"context"
	"errors"
	"fmt"
	"glitchz/pkg/models"
	"glitchz/pkg/services/password"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthService interface {
	CreateUser(data models.User) error
	UpdateUser(filter bson.D, update bson.D) (mongo.UpdateResult, error)
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
	id := primitive.NewObjectID()
	hashedPass, err := password.HashPassword(data.Password)

	if err != nil {
		return err
	}

	if user, _ := GetUserByEmail(a, data.Email); user.Email == data.Email {
		err = errors.New("this user already exists")
		return err
	}

	newUser := models.User{
		ID:        id,
		Username:  fmt.Sprintf("@" + data.Username),
		Email:     data.Email,
		Password:  hashedPass,
		CreatedAT: time.Now(),
	}

	_, err = a.col.InsertOne(a.ctx, newUser)

	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(a *authService, email string) (models.User, error) {
	user := models.User{}
	filter := bson.D{{Key: "email", Value: email}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *authService) UpdateUser(filter bson.D, update bson.D) (mongo.UpdateResult, error) {
	result, err := a.col.UpdateOne(a.ctx, filter, update, options.Update())
	if err != nil {
		return mongo.UpdateResult{}, err
	}
	return *result, nil
}
