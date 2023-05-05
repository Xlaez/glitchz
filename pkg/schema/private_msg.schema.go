package schema

type SendMsgReq struct {
	ContactID string `form:"contactId" binding:"required"`
	Text      string `form:"text"`
}

type DeleteMsgReq struct {
	ID string `uri:"id" binidng:"required"`
}

type ReadMsgReq struct {
	ID string `uri:"id" binidng:"required"`
}

type AddReactionReq struct {
	MsgID    string `json:"msgId" binding:"required"`
	Reaction string `json:"reaction" binding:"required"`
}
