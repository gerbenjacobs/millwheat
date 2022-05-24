package millwheat

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserEmailUniqueness = errors.New("email address already in use")

	ErrEmailNotFound = errors.New("no user found for email")
	ErrWrongPassword = errors.New("password doesn't match")

	ErrPageNotFound     = errors.New("page not found")
	ErrMethodNotAllowed = errors.New("method not allowed")

	ErrItemNotFound          = errors.New("item not found")
	ErrItemNotEnoughQuantity = errors.New("not enough quantity")
	ErrNoItems               = errors.New("no items")
)
