package schema

import (
	"glitchz/pkg/models"
	"glitchz/pkg/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,min=6"`
	Username string `json:"username" binding:"required,alphanum,min=2"`
}

type AddUserRes struct {
	Code string `json:"code"`
}

type VerifyEmailReq struct {
	Code string `json:"code" binding:"required,min=6,max=6"`
}

type UserRes struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"id"`
	Username    string             `json:"username,omitempty" bson:"username"`
	Email       string             `json:"email,omitempty" bson:"email"`
	Pics        string             `json:"pics,omitempty" bson:"pics"`
	Role        string             `json:"role" bson:"role" default:"user"`
	Bio         string             `json:"bio,omitempty" bson:"bio"`
	LinkedIN    string             `json:"linkedIn,omitempty" bson:"linkedIn"`
	Github      string             `json:"github,omitempty" bson:"github"`
	Twitter     string             `json:"twitter,omitempty" bson:"twitter"`
	IsSuspended bool               `json:"isSuspended" bson:"isSuspended" default:"false"`
	NBContacts  int                `json:"nbContacts" bson:"nbContacts" default:"0"`
	NBReviews   int                `json:"nbReviews" bson:"nbReviews" default:"0"`
}

type LoginReq struct {
	EmailOrUsername string `json:"emailOrUsername" binding:"required"`
	Password        string `json:"password" binding:"required,alphanum,min=6"`
}

type LoginRes struct {
	Token *services.TokenService `json:"token"`
	User  models.User            `json:"user"`
}

type RefreshAccessTokenReq struct {
	Token string `json:"token" binding:"required"`
}

type RequestPasswordResetReq struct {
	Email string `json:"email" binding:"email,email"`
}

type ConfirmPassword struct {
	Password string `json:"password" binding:"required,min=6,alphanum"`
}

type ResetPasswordReq struct {
	Code     string `json:"code" binding:"required,min=6,max=6"`
	Password string `json:"password" binding:"required,min=6,alphanum"`
}

type GetUserByIdReq struct {
	UserID string `uri:"userId" binding:"required"`
}

type GetUserByUsername struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

type UpdateProfileReq struct {
	Bio      string `json:"bio"`
	LinkedIN string `json:"linkedIn"`
	Github   string `json:"github"`
	Twitter  string `json:"twitter"`
}

type GetUsersReq struct {
	Limit int `form:"limit" binding:"required,min=5"`
	Page  int `form:"page" binding:"required,min=1"`
}
