package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"id"`
	User1      primitive.ObjectID   `json:"user1,omitempty" bson:"user1"`
	User2      primitive.ObjectID   `json:"user2,omitempty" bson:"user2"`
	Pending    bool                 `json:"pending" bson:"pending" default:"false"`
	BlockedIDs []primitive.ObjectID `json:"blockedIds" bson:"blockedIds"`
	SentAT     time.Time            `json:"sentAt" bson:"sentAt"`
	AcceptedAT time.Time            `json:"acceptedAt" bson:"acceptedAt"`
}
