package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostService interface {
	NewPost()
}

type postService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewPostService(col *mongo.Collection, ctx context.Context) PostService {
	return &postService{
		col: col,
		ctx: ctx,
	}
}

func (p *postService) NewPost() {}
