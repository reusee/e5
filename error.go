package e5

import (
	"errors"
	"strings"
)

// Error represents a chain of errors
type Error []error

// Is reports whether any error in the chain matches target
func (c Error) Is(target error) bool {
	for _, err := range c {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As reports whether any error in the chain matches target.
// And if so, assign the first matching error to target
func (c Error) As(target interface{}) bool {
	for _, err := range c {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

// Error implements error interface
func (c Error) Error() string {
	var b strings.Builder
	for i, err := range c {
		var str string
		if leveled, ok := err.(ErrorLeveler); ok {
			if leveled.ErrorLevel() <= ErrorLevel {
				str = err.Error()
			}
		} else {
			str = err.Error()
		}
		if i > 0 && len(str) > 0 && b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(str)
	}
	return b.String()
}

// Chain chains two errors
func Chain(err error, prev error) Error {
	chain, ok := err.(Error)
	if ok {
		chain = append(chain, prev)
		return chain
	}
	chain, ok = prev.(Error)
	if ok {
		chain = append(chain, err)
		return chain
	}
	return Error{err, prev}
}

// With returns a WrapFunc that wraps an error value
func With(err error) WrapFunc {
	return func(prev error) error {
		if prev == nil {
			return nil
		}
		return Chain(err, prev)
	}
}
