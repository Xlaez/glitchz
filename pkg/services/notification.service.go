package services

import (
	"context"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationService interface {
	NewNotification(data models.Notification) (*mongo.InsertOneResult, error)
	InsertMany(data []interface{}) error
}

type notificationService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewNotificationService(col *mongo.Collection, ctx context.Context) NotificationService {
	return &notificationService{
		col: col,
		ctx: ctx,
	}
}

func (n *notificationService) NewNotification(data models.Notification) (*mongo.InsertOneResult, error) {
	result, err := n.col.InsertOne(n.ctx, data)

	if err != nil {
		return &mongo.InsertOneResult{}, err
	}
	return result, nil
}

func (n *notificationService) InsertMany(data []interface{}) error {
	if _, err := n.col.InsertMany(n.ctx, []interface{}{}); err != nil {
		return err
	}

	return nil
}
