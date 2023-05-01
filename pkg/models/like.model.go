package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Likes struct {
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
}
