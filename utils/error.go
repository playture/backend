package utils

import (
	"errors"
)

func WrapError(msg string, err ...error) error {
	if len(err) == 0 {
		return errors.New(msg)
	}
	wrappedErr := errors.New(msg)
	return errors.Join(wrappedErr, errors.Join(err...))
}
