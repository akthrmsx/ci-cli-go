package main

import (
	"errors"
	"fmt"
)

var ErrValidation = errors.New("validation failed")

type stepError struct {
	step    string
	message string
	cause   error
}

func (x *stepError) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", x.step, x.message, x.cause)
}

func (x *stepError) Is(target error) bool {
	y, ok := target.(*stepError)

	if !ok {
		return false
	}

	return x.step == y.step
}

func (x *stepError) Unwrap() error {
	return x.cause
}
