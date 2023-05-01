package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadBy struct {
	UserID primitive.ObjectID `bson:"userId"`
	ReadAT time.Time          `bson:"readAt"`
}
