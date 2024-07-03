package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	SecretKey string
}

// createToken implements Maker.
func (j *JWTMaker) CreateToken(username string, role string, duration time.Duration) (token string, payload *Payload, err error) {
	payload, err = NewPayload(duration, username, role)
	if err != nil {
		return "nil", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err = jwtToken.SignedString([]byte(j.SecretKey))

	fmt.Println("the token generated is", token)

	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be atleast %d characters", minSecretKeySize)
	}
	return &JWTMaker{SecretKey: secretKey}, nil
}
