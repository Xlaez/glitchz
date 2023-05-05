package services

import (
	"context"
	"errors"
	"glitchz/pkg/models"
	"glitchz/pkg/services/token"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tokens struct {
	AccessToken  string
	RefreshToken string
}
type TokenService interface {
	GenerateToken(userId primitive.ObjectID, duration time.Duration) (*tokens, error)
	DeleteToken(token string) error
}

type tokenService struct {
	col *mongo.Collection
	t   token.Maker
	ctx context.Context
}

func NewTokenService(col *mongo.Collection, t token.Maker, ctx context.Context) TokenService {
	return &tokenService{
		col: col,
		t:   t,
		ctx: ctx,
	}
}

func (t *tokenService) GenerateToken(userId primitive.ObjectID, duration time.Duration) (*tokens, error) {

	access_token, err := t.t.CreateToken(userId, duration)

	if err != nil {
		return &tokens{}, err
	}

	refresh_token, err := t.t.CreateToken(userId, 10000000*time.Second)

	if err != nil {
		return &tokens{}, err
	}

	_, err = t.col.InsertOne(t.ctx, models.Token{
		ID:        primitive.NewObjectID(),
		Token:     refresh_token,
		UserID:    userId,
		Type:      "refresh",
		ExpiresAT: 100000 * time.Second,
		CreatedAT: time.Now(),
	})

	if err != nil {
		return &tokens{}, err
	}

	return &tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil

}

func (t *tokenService) DeleteToken(token string) error {
	filter := bson.D{{Key: "token", Value: token}}
	result, err := t.col.DeleteOne(t.ctx, filter, options.Delete())
	if err != nil || result.DeletedCount == 0 {
		return errors.New("cannot delete token, token not found")
	}
	return nil
}
