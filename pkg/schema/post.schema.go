package schema

type NewPostReq struct {
	Text   string `json:"text" binding:"required"`
	Public bool   `json:"public" binding:"required"`
}
