package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserID     primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
	Text       string             `json:"text,omitempty" bson:"text"`
	FIle       string             `json:"file,omitempty" bson:"file"`
	Likes      []Likes            `json:"likes" bson:"likes"`
	NBLikes    int                `json:"nbLikes" bson:"nbLikes" default:"0"`
	NBComments int                `json:"nbComments" bson:"nbComments" default:"0"`
	Public     bool               `json:"public" bson:"public" default:"true"`
	CreatedAT  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAT  time.Time          `json:"updatedAt" bson:"updatedAt"`
}
