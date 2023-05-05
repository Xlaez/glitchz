package services

import (
	"context"
	"fmt"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContactService interface {
	SendReq(data models.Contact) (string, error)
	UpdateReq(filter bson.D, update bson.D) (models.Contact, error)
	GetContacts(bson.D, *options.FindOptions) ([]models.Contact, int64, error)
}

type contactService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewContactService(col *mongo.Collection, ctx context.Context) ContactService {
	return &contactService{
		col: col,
		ctx: ctx,
	}
}

func (c *contactService) SendReq(data models.Contact) (string, error) {
	result, err := c.col.InsertOne(c.ctx, data, options.InsertOne())

	if err != nil {
		return "", err
	}

	return fmt.Sprint(result), nil
}

func (c *contactService) UpdateReq(filter, update bson.D) (models.Contact, error) {
	result := models.Contact{}
	c.col.FindOneAndUpdate(c.ctx, filter, update, options.FindOneAndUpdate()).Decode(&result)
	return result, nil
}

func (c *contactService) GetContacts(filter bson.D, options *options.FindOptions) ([]models.Contact, int64, error) {
	contacts := []models.Contact{}
	cursor, err := c.col.Find(c.ctx, filter, options)

	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(c.ctx, &contacts); err != nil {
		return nil, 0, err
	}

	count, err := c.col.CountDocuments(c.ctx, filter)

	if err != nil {
		return nil, 0, err
	}

	return contacts, count, nil
}

func (c *contactService) DeleteContact()
