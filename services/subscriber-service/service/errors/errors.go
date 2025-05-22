package errors

import "errors"

var (
	ErrEmailAlreadySubscribed = errors.New("email already subscribed to the newsletter")
	ErrEmailNotSubscribed     = errors.New("email not subscribed to the newsletter")
)
