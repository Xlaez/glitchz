package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid token key size, must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(user_id primitive.ObjectID, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(user_id, duration)

	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	var err error
	if err = maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil); err != nil {
		return nil, err
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
