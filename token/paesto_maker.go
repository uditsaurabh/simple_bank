package token

import (
	"errors"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

const KeySize = 32

var (
	errBadKeySize = errors.New("chacha20poly1305: bad key length")
	errAuthFailed = errors.New("chacha20poly1305: message authentication has failed")
)

type PasetoMaker struct {
	Paseto       *paseto.V2
	SymmetricKey []byte
}

// CreateToken implements Maker.
func (p *PasetoMaker) CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(duration, username, role)
	if err != nil {
		return "", payload, err
	}
	paestoToken, err := p.Paseto.Encrypt(p.SymmetricKey, payload, nil)
	if err != nil {
		return "", payload, err
	}
	return paestoToken, payload, nil
}

// VerifyToken implements Maker.
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.Paseto.Decrypt(token, p.SymmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, errAuthFailed
	}
	return payload, nil
}

func NewPasetoMaker(symmetrickey string) (Maker, error) {
	if len(symmetrickey) != chacha20poly1305.KeySize {
		return nil, errBadKeySize
	}
	maker := &PasetoMaker{
		Paseto:       paseto.NewV2(),
		SymmetricKey: []byte(symmetrickey),
	}
	return maker, nil
}
