package schema

type NewCommentReq struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parentId"`
	PostID   string `json:"postId" binding:"required"`
}
