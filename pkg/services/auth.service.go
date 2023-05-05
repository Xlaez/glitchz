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
	LoginUser(emailOrPassword string, password string) (models.User, error)
	GetUserById(userId primitive.ObjectID) (models.User, error)
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

	users, _ := GetUsers(a, bson.D{{Key: "username", Value: data.Username}})
	if users != nil {
		return errors.New("username has been taken, try another")
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

func (a *authService) LoginUser(emailOrUsername string, pass string) (models.User, error) {

	value := make([]primitive.D, 0)
	value = append(value, bson.D{{Key: "email", Value: emailOrUsername}})
	value = append(value, bson.D{{Key: "username", Value: fmt.Sprint("@", emailOrUsername)}})

	filter := bson.D{{Key: "$or", Value: value}}

	var user models.User

	if err := a.col.FindOne(a.ctx, filter, options.FindOne()).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, errors.New("username or email wrong")
		}
		return models.User{}, err
	}

	if err := password.ComparePassword(pass, user.Password); err != nil {
		return models.User{}, errors.New("password does not match")
	}

	return user, nil
}

func GetUserByEmail(a *authService, email string) (models.User, error) {
	user := models.User{}
	filter := bson.D{{Key: "email", Value: email}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *authService) GetUserById(userId primitive.ObjectID) (models.User, error) {
	user := models.User{}
	filter := bson.D{primitive.E{Key: "id", Value: userId}}

	if err := a.col.FindOne(a.ctx, filter).Decode(&user); err == mongo.ErrNoDocuments && err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUsers(a *authService, filter bson.D) ([]models.User, error) {
	cursor, err := a.col.Find(a.ctx, filter, options.Find().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err = cursor.All(a.ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (a *authService) UpdateUser(filter bson.D, update bson.D) (mongo.UpdateResult, error) {
	result, err := a.col.UpdateOne(a.ctx, filter, update, options.Update())
	if err != nil {
		return mongo.UpdateResult{}, err
	}
	return *result, nil
}
