package view

import "errors"

var (
	ErrExpiredPlan        = errors.New("your plan expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTooBigLength       = errors.New("too big length")
)
