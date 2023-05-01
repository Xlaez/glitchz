package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommentService interface {
	NewComment()
}

type commentService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewCommentService(col *mongo.Collection, ctx context.Context) CommentService {
	return &commentService{
		col: col,
		ctx: ctx,
	}
}

func (c *commentService) NewComment() {}
