package services

import (
	"context"
	"glitchz/pkg/models"
	"glitchz/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	GetUser(filter bson.D) (models.User, error)
	UpdateUser(filter bson.D, update bson.D) (*mongo.UpdateResult, error)
	GetUsers(filter bson.D, options *options.FindOptions) ([]schema.UserRes, int64, error)
	UpdateMany(update mongo.UpdateManyModel) error
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

func (u *userService) GetUser(filter bson.D) (models.User, error) {
	var user models.User
	if err := u.col.FindOne(u.ctx, filter, options.FindOne()).Decode(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userService) UpdateUser(filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	result, err := u.col.UpdateOne(u.ctx, filter, update, options.Update())
	if err != nil {
		return &mongo.UpdateResult{}, err
	}
	return result, nil
}

func (u *userService) GetUsers(filter bson.D, options *options.FindOptions) ([]schema.UserRes, int64, error) {
	users := []schema.UserRes{}
	cursor, err := u.col.Find(u.ctx, filter, options, options.SetAllowDiskUse(true))

	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(u.ctx, &users); err != nil {
		return nil, 0, err
	}

	count, err := u.col.CountDocuments(u.ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (u *userService) UpdateMany(update mongo.UpdateManyModel) error {
	_, err := u.col.UpdateMany(u.ctx, update, options.Update())

	if err != nil {
		return err
	}

	return nil
}
