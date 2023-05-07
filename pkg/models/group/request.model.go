package group

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	GroupID primitive.ObjectID `json:"groupId,omitempty" bson:"groupId"`
	Msg     string             `json:"msg,omitempty" bson:"msg" required:"false"`
	UserID  primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
	SentAT  time.Time          `json:"sentAt" bson:"sentAt"`
}
