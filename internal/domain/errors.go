package domain

import "errors"

var (
	ErrNoSuchEntity = errors.New("no such entity")
	ErrInvalidInput = errors.New("invalid input value")
)
