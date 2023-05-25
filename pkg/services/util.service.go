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

func GetUserContacts(ctx context.Context, col *mongo.Collection, userId primitive.ObjectID) ([]models.Contact, error) {
	contacts := []models.Contact{}
	filter := make([]bson.D, 0)
	filter = append(filter, bson.D{{Key: "user1", Value: userId}, {Key: "pending", Value: false}})
	filter = append(filter, bson.D{{Key: "user2", Value: userId}, {Key: "pending", Value: false}})
	cursor, err := col.Find(ctx, bson.D{{Key: "$or", Value: filter}}, options.Find())

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}

	if err = cursor.Close(ctx); err != nil {
		return nil, err
	}
	fmt.Print(contacts)
	return contacts, nil
}
