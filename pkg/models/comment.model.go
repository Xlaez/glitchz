package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
	Content   string             `json:"content" bson:"content"`
	NBReply   int                `json:"nbReplies" bson:"nbReplies" default:"0"`
	NBLikes   int                `json:"nbLikes" bson:"nbLikes" default:"0"`
	Likes     []Likes            `json:"likes" bson:"likes"`
	ParentID  primitive.ObjectID `json:"parentId,omitempty" bson:"parentId"`
	PostID    primitive.ObjectID `json:"postId" bson:"postId" `
	CreatedAT time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAT time.Time          `json:"updatedAt" bson:"updatedAt"`
}
