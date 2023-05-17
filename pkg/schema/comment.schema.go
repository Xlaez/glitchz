package schema

type NewCommentReq struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parentId"`
	PostID   string `json:"postId" binding:"required"`
}

type GetCommentsReq struct {
	Limit  int64  `form:"limit" binding:"required,min=5"`
	Page   int64  `form:"page" binding:"required,min=1"`
	PostID string `form:"postId" binding:"required"`
}

type GetCommentsRepliesReq struct {
	Limit     int64  `form:"limit" binding:"required,min=5"`
	Page      int64  `form:"page" binding:"required,min=1"`
	CommentID string `form:"commentId" binding:"required"`
}

type UpdateCommentReq struct {
	CommentID string `json:"commentId" binding:"required"`
	Text      string `json:"text" binding:"required"`
}

type DeleteCommentReq struct {
	ID string `uri:"id" binding:"required"`
}
