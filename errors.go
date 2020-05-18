package svc

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserEmailUniqueness = errors.New("email address already in use")

	ErrCharacterNotFound = errors.New("character not found")

	ErrEmailNotFound = errors.New("no user found for email")
	ErrWrongPassword = errors.New("password doesn't match")

	ErrPageNotFound     = errors.New("page not found")
	ErrMethodNotAllowed = errors.New("method not allowed")
)
