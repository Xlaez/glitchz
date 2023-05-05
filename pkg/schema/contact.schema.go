package schema

type SendReq struct {
	UserID string `json:"userId" binding:"required"`
}

type AcceptReq struct {
	ID string `json:"Id" binding:"required"`
}

type BlockReq struct {
	UserID string `json:"userId" binding:"required"`
	ID     string `json:"contactId" binding:"required"`
}

type GetContactByUser struct {
	UserID string `form:"userId" binding:"required"`
}

type GetContactByID struct {
	ID string `form:"id" binidng:"required"`
}
