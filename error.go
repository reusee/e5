package e5

import (
	"strings"
)

// Error represents a chain of errors
type Error []error

// Unwrap returns wrapped errors
func (c Error) Unwrap() []error {
	return c
}

// Error implements error interface
func (c Error) Error() string {
	var b strings.Builder
	for i, err := range c {
		str := err.Error()
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
