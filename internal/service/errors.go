package service

import "errors"

var (
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrEmailAlreadyExist  = errors.New("email already exist")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
