package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Msg struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	ContactID primitive.ObjectID `json:"contactId,omitempty" bson:"contactId"`
	Sender    primitive.ObjectID `json:"sender,omitempty" bson:"sender"`
	ISDeleted bool               `json:"isDeleted" bson:"isDeleted" default:"false"`
	ReadBY    []ReadBy           `bson:"readBy"`
	Reaction  []Reaction         `bson:"reaction"`
	Msg       Message            `json:"msg" bson:"msg"`
	SentAT    time.Time          `json:"sentAt" bson:"sentAt"`
}

type Message struct {
	Text string `json:"text,omitempty" bson:"text"`
	File string `json:"file,omitempty" bson:"file"`
}

type Reaction struct {
	UserID   primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
	Reaction string             `json:"reaction"` // enum["like", "laugh", "dislike", "angry", "sad", "clap"]
}
