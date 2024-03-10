package util

import "errors"

var (
	ErrUnknown         = errors.New("Invalid login or password")
	ErrInvalidArgument = errors.New("invalid argument passed")
)
