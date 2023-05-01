package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"id"`
	Username    string             `json:"username,omitempty" bson:"username"`
	Email       string             `json:"email,omitempty" bson:"email"`
	Password    string             `json:"password,omitempty" bson:"password"`
	Pics        string             `json:"pics,omitempty" bson:"pics"`
	Role        string             `json:"role" bson:"role" default:"user"`
	Bio         string             `json:"bio,omitempty" bson:"bio"`
	LinkedIN    string             `json:"linkedIn,omitempty" bson:"linkedIn"`
	Github      string             `json:"github,omitempty" bson:"github"`
	Twitter     string             `json:"twitter,omitempty" bson:"twitter"`
	IsVerified  bool               `json:"isVerified" bson:"isVerified" default:"false"`
	IsSuspended bool               `json:"isSuspended" bson:"isSuspended" default:"false"`
	NBContacts  int                `json:"nbContacts" bson:"nbContacts" default:"0"`
	NBReviews   int                `json:"nbReviews" bson:"nbReviews" default:"0"`
	CreatedAT   time.Time          `json:"createdAt" bson:"createdAt"`
}
