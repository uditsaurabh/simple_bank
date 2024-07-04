package token

import "time"

type TokenCreate interface {
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)
}

type TokenVerify interface {
	VerifyToken(token string) (*Payload, error)
}

type Maker interface {
	TokenCreate
	TokenVerify
}
