package group

import (
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupMsg struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"id"`
	GroupID   primitive.ObjectID `json:"groupId,omitempty" bson:"groupId"`
	Sender    primitive.ObjectID `json:"sender,omitempty" bson:"sender"`
	ISDeleted bool               `json:"isDeleted" bson:"isDeleted" default:"false"`
	ReadBY    []models.ReadBy    `json:"readBy" bson:"readBy"`
	Msg       models.Message     `json:"msg" bson:"msg"`
}
