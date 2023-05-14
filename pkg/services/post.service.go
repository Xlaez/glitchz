package services

import (
	"context"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostService interface {
	NewPost(insertData models.Post) error
	GetPost(filter bson.D) (models.Post, error)
	GetPosts(filter bson.D, options *options.FindOptions) (int64, []models.Post, error)
	Update(filter bson.D, updateObj bson.D) (*mongo.UpdateResult, error)
	Delete(filter bson.D) error
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

func (p *postService) GetPost(filter bson.D) (models.Post, error) {
	post := models.Post{}
	if err := p.col.FindOne(p.ctx, filter).Decode(&post); err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (p *postService) GetPosts(filter bson.D, options *options.FindOptions) (int64, []models.Post, error) {
	posts := []models.Post{}
	cursor, err := p.col.Find(p.ctx, filter, options)
	if err != nil {
		return 0, nil, err
	}

	if err = cursor.All(p.ctx, &posts); err != nil {
		return 0, nil, err
	}
	if err = cursor.Close(p.ctx); err != nil {
		return 0, nil, err
	}

	count, err := p.col.CountDocuments(p.ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	return count, posts, nil
}

func (p *postService) Update(filter bson.D, updateObj bson.D) (*mongo.UpdateResult, error) {
	result, err := p.col.UpdateOne(p.ctx, filter, updateObj, options.Update())
	if err != nil {
		return &mongo.UpdateResult{}, err
	}
	return result, nil
}

func (p *postService) Delete(filter bson.D) error {
	if _, err := p.col.DeleteOne(p.ctx, filter, options.Delete()); err != nil {
		return err
	}
	return nil
}
