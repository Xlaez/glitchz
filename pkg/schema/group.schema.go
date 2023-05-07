package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type NewGroupReq struct {
	Public  bool     `json:"public"`
	Name    string   `json:"name" binding:"required"`
	Members []string `json:"members" binding:"required"`
}

type GetGroupByIDReq struct {
	ID string `uri:"id" binding:"required"`
}

type GetGroupsReq struct {
	Limit int64 `form:"limit" binding:"required,min=5"`
	Page  int64 `form:"page" binidng:"required,min=1"`
}

type AddMembersToGroup struct {
	IDs     []primitive.ObjectID `json:"ids" binding:"required"`
	GroupID string               `json:"groupId" binding:"required"`
}

type JoinGroup struct {
	GroupID string `uri:"groupId" binding:"required"`
}
