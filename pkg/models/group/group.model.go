package group

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Pics       string             `json:"pics,omitempty" bson:"pics" default:"https://unsplash.com/bot1-img.png"`
	Name       string             `json:"name" bson:"name"`
	Public     bool               `json:"public" bson:"public"`
	NBRequests int                `json:"nbRequests" bson:"nbRequests"`
	BlockedIDS BlockedIDs         `json:"blockedIds" bson:"blockedIds"`
	Members    []Members          `json:"members" bson:"members"`
	Admins     []Members          `json:"admins" bson:"admins"`
	CreatedAT  time.Time          `json:"createdAt" bson:"createdAt"`
}

type BlockedIDs struct {
	UserID primitive.ObjectID `json:"userId" bson:"userId"`
}

type Members struct {
	UserID   primitive.ObjectID `json:"userId" bson:"userId"`
	JoinedAT time.Time          `json:"joinedAt" bson:"joinedAt"`
}
