package schema

type NewPostReq struct {
	Text   string `form:"text" binding:"required"`
	Public bool   `form:"public" `
}

type GetPostByIDReq struct {
	ID string `uri:"id" binding:"required"`
}

type GetPostsByUserIDReq struct {
	Limit  int64  `form:"limit" binding:"required,min=5"`
	Page   int64  `form:"page" binding:"required,min=1"`
	UserID string `form:"userId" binding:"required"`
}

type GetAllPostsReq struct {
	Limit int64 `form:"limit" binding:"required,min=5"`
	Page  int64 `form:"page" binding:"required,min=1"`
}

type UpdatePostReq struct {
	Text   string `json:"text"`
	ID     string `json:"id" binding:"required"`
	Public string `json:"public"`
}

type DeletePostReq struct {
	ID string `uri:"id" binding:"required"`
}
