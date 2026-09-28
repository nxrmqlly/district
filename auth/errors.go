package auth

import "errors"

var (
	ErrUsernameTaken = errors.New("username already taken")
	ErrEmailInUse    = errors.New("email already in use")
)
