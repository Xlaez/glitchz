package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId,omitempty"`
	Msg       string             `bson:"msg" json:"msg"`
	Image     string             `bson:"img" json:"img"`
	Seen      bool               `bson:"seen" json:"seen" default:"false"`
	CreatedAT time.Time          `bson:"createdAt" json:"createdAt"`
}
