package services

import (
	"context"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostService interface {
	NewPost(insertData models.Post) error
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

func (p *postService) NewPost(insertData models.Post) error {
	if _, err := p.col.InsertOne(p.ctx, insertData, options.InsertOne()); err != nil {
		return err
	}
	return nil
}
