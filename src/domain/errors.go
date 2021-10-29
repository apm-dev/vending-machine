package domain

import "github.com/pkg/errors"

var (
	ErrInternalServer    = errors.New("internal server error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidToken      = errors.New("unauthorized user")
	ErrUnauthorized      = errors.New("unauthorized user")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrInvalidCoin       = errors.New("invalid coin, use 5, 10, 20, 50, 100 cent coins")
)
