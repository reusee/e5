package e5

import (
	"errors"
	"strings"
)

// Error represents multiple errors
type Error []error

func (c Error) copy() Error {
	ret := make(Error, len(c))
	copy(ret, c)
	return ret
}

// Is reports whether any error in the slice matches target
func (c Error) Is(target error) bool {
	for _, err := range c {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As reports whether any error in the slice matches target.
// And if so, assign the first matching error to target
func (c Error) As(target interface{}) bool {
	for _, err := range c {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

// Unwrap returns all errors
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

// Join joins two errors
func Join(err error, prev error) Error {
	e1, ok1 := err.(Error)
	e2, ok2 := prev.(Error)
	if ok1 && !ok2 {
		if e1.Is(prev) {
			return e1
		}
		return append(e1.copy(), prev)
	}
	if !ok1 && ok2 {
		if e2.Is(err) {
			return e2
		}
		return append(e2.copy(), err)
	}
	if !ok1 && !ok2 {
		if errors.Is(err, prev) {
			return Error{err}
		}
		return Error{err, prev}
	}
	e1 = e1.copy()
	for _, e := range e2 {
		if e1.Is(e) {
			continue
		}
		e1 = append(e1, e)
	}
	return e1
}

// With returns a WrapFunc that wraps an error value
func With(err error) WrapFunc {
	return func(prev error) error {
		if prev == nil {
			return nil
		}
		return Join(err, prev)
	}
}
