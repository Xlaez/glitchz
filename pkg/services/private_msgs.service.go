package services

import (
	"context"
	"glitchz/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PrivateMsgs interface {
	SendMsg(models.Msg) error
	Updatemsg(filter, update bson.D) error
	GetMsgByID(id primitive.ObjectID) (models.Msg, error)
}

type privateMsgs struct {
	col *mongo.Collection
	ctx context.Context
}

func NewPrivateMsgs(col *mongo.Collection, ctx context.Context) PrivateMsgs {
	return &privateMsgs{
		col: col,
		ctx: ctx,
	}
}

func (p *privateMsgs) SendMsg(data models.Msg) error {
	if _, err := p.col.InsertOne(p.ctx, data, options.InsertOne()); err != nil {
		return err
	}
	return nil
}

func (p *privateMsgs) Updatemsg(filter, update bson.D) error {
	if _, err := p.col.UpdateOne(p.ctx, filter, update, options.Update()); err != nil {
		return err
	}
	return nil
}

func (p *privateMsgs) GetMsgByID(id primitive.ObjectID) (models.Msg, error) {
	msg := models.Msg{}

	if err := p.col.FindOne(p.ctx, bson.D{primitive.E{Key: "_id", Value: id}}, options.FindOne()).Decode(&msg); err != nil {
		return models.Msg{}, err
	}

	return msg, nil
}

func (p *privateMsgs) GetRecentMsgs(contactIds []primitive.ObjectID, userId primitive.ObjectID) {

	// contacts, err := p.col.Aggregate(p.ctx, )
}
