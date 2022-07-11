package token

import (
	"errors"
	"time"
)

var ErrInvalidToken = errors.New("token is invalid")
var ErrExpiredToken = errors.New("token has expired")

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}