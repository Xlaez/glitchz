package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"id"`
	Token     string             `json:"token,omitempty" bson:"token"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Type      string             `json:"type" bson:"type"` // ["refresh", "access"]
	ExpiresAT time.Time          `json:"expiresAt" bson:"expiresAt"`
	CreatedAT time.Time          `json:"createdAt" bson:"createdAt"`
}
