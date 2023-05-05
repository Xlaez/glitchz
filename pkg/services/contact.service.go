package services

import (
	"context"
	"fmt"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContactService interface {
	GetContact(bson.D, *options.FindOneOptions) (models.Contact, error)
	SendReq(data models.Contact) (string, error)
	UpdateReq(filter bson.D, update bson.D) (models.Contact, error)
	GetContacts(bson.D, *options.FindOptions) ([]models.Contact, int64, error)
	GetConvsByUserId(id primitive.ObjectID) ([]models.Contact, error)
	DeleteContact(filter bson.D) (models.Contact, error)
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

func (c *contactService) DeleteContact(filter bson.D) (models.Contact, error) {
	contact := models.Contact{}
	err := c.col.FindOneAndDelete(c.ctx, filter, options.FindOneAndDelete()).Decode(&contact)
	if err != nil {
		return models.Contact{}, err
	}

	return contact, nil
}

func (c *contactService) GetContact(filter bson.D, options *options.FindOneOptions) (models.Contact, error) {
	contacts := models.Contact{}
	if err := c.col.FindOne(c.ctx, filter, options).Decode(&contacts); err != nil {
		return models.Contact{}, err
	}
	return contacts, nil
}

func (c *contactService) GetConvsByUserId(id primitive.ObjectID) ([]models.Contact, error) {
	contacts := []models.Contact{}

	filter := make([]bson.D, 0)
	filter = append(filter, bson.D{{Key: "user1", Value: id}, {Key: "pending", Value: false}})
	filter = append(filter, bson.D{{Key: "user2", Value: id}, {Key: "pending", Value: false}})
	curosr, err := c.col.Find(c.ctx, bson.D{{Key: "$or", Value: filter}}, options.Find())
	if err != nil {
		return nil, err
	}
	if err = curosr.All(c.ctx, &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}
