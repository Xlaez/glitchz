package token

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Maker interface {
	CreateToken(user_id primitive.ObjectID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
