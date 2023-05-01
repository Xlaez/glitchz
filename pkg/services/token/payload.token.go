package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID          `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	IssuedAt  time.Time          `json:"issued_at"`
	ExpiresAt time.Time          `json:"expires_at"`
}

func NewPayLoad(user_id primitive.ObjectID, duration time.Duration) (*Payload, error) {
	token_id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        token_id,
		UserID:    user_id,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
