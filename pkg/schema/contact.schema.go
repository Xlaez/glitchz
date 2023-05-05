package schema

type SendReq struct {
	UserID string `json:"userId" binding:"required"`
}

type AcceptReq struct {
	ID string `json:"Id" binding:"required"`
}
