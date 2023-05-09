package schema

type NewGroupMsgReq struct {
	GroupID string `form:"groupId" binding:"required"`
	Text    string `form:"text" binding:"required"`
}
