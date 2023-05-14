package services

import (
	"context"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentService interface {
	NewComment(data models.Comment) error
	GetComments(filter bson.D, options *options.FindOptions) ([]models.Comment, error)
	UpdateComment(filter bson.D, update bson.D) (*models.Comment, error)
	GetComment(filter bson.D) (*models.Comment, error)
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

func (c *commentService) NewComment(data models.Comment) error {
	_, err := c.col.InsertOne(c.ctx, data, options.InsertOne())
	if err != nil {
		return err
	}
	return nil
}

func (c *commentService) GetComments(filter bson.D, options *options.FindOptions) ([]models.Comment, error) {
	comments := []models.Comment{}
	cursor, err := c.col.Find(c.ctx, filter, options)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(c.ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *commentService) UpdateComment(filter, update bson.D) (*models.Comment, error) {
	result := models.Comment{}
	if err := c.col.FindOneAndUpdate(c.ctx, filter, update, options.FindOneAndUpdate()).Decode(&result); err != nil {
		return &result, err
	}
	return &result, nil
}

func (c *commentService) GetComment(filter bson.D) (*models.Comment, error) {
	comment := models.Comment{}
	if err := c.col.FindOne(c.ctx, filter, options.FindOne()).Decode(&comment); err != nil {
		return &models.Comment{}, err
	}
	return &comment, nil
}
