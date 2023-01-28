package e5

import (
	"errors"
	"strings"
)

// Error represents multiple errors
type Error map[error]struct{}

func (c Error) copy() Error {
	ret := make(Error, len(c))
	for e := range c {
		ret[e] = struct{}{}
	}
	return ret
}

// Is reports whether any error in the slice matches target
func (c Error) Is(target error) bool {
	if _, ok := c[target]; ok {
		return true
	}
	for err := range c {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As reports whether any error in the slice matches target.
// And if so, assign the first matching error to target
func (c Error) As(target interface{}) bool {
	for err := range c {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

// Unwrap returns all errors
func (c Error) Unwrap() []error {
	ret := make([]error, 0, len(c))
	for err := range c {
		ret = append(ret, err)
	}
	return ret
}

// Error implements error interface
func (c Error) Error() string {
	var b strings.Builder
	i := 0
	for err := range c {
		str := err.Error()
		if i > 0 && len(str) > 0 && b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(str)
		i++
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
		ret := e1.copy()
		ret[prev] = struct{}{}
		return ret
	}
	if !ok1 && ok2 {
		if e2.Is(err) {
			return e2
		}
		ret := e2.copy()
		ret[err] = struct{}{}
		return ret
	}
	if !ok1 && !ok2 {
		if errors.Is(err, prev) {
			return Error{
				err: struct{}{},
			}
		}
		return Error{
			err:  struct{}{},
			prev: struct{}{},
		}
	}
	e1 = e1.copy()
	for e := range e2 {
		if e1.Is(e) {
			continue
		}
		e1[e] = struct{}{}
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
