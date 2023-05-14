package schema

type NewPostReq struct {
	Text   string `form:"text" binding:"required"`
	Public bool   `form:"public" binding:"required"`
}

type GetPostByIDReq struct {
	ID string `uri:"id" binding:"required"`
}

type GetPostsByUserIDReq struct {
	Limit  int64  `form:"limit" binding:"required,min=5"`
	Page   int64  `form:"page" binding:"required,min=1"`
	UserID string `form:"userId" binding:"required"`
}
