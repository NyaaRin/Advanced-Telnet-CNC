package database

import "errors"

var (
	ErrCredentialsInvalid = errors.New("invalid credentials were supplied")
	ErrKnownUser          = errors.New("user exists")
	ErrUnknownUser        = errors.New("user does not exist")
)
