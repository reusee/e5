package e5

import (
	"fmt"
	"sync"
)

// ErrInfo represents a lazy-evaluaed formatted string
type ErrInfo struct {
	format     string
	str        string
	args       []any
	formatOnce sync.Once
}

var _ error = new(ErrInfo)

// Error implements error interface
func (i *ErrInfo) Error() string {
	i.formatOnce.Do(func() {
		i.str = fmt.Sprintf(i.format, i.args...)
	})
	return i.str
}

// Info returns a WrapFunc that wraps an *ErrInfo error value
func Info(format string, args ...any) WrapFunc {
	return With(&ErrInfo{
		format: format,
		args:   args,
	})
}

func (w WrapFunc) WithInfo(format string, args ...any) WrapFunc {
	return w.With(Info(format, args...))
}

func (c CheckFunc) WithInfo(format string, args ...any) CheckFunc {
	return c.With(Info(format, args...))
}
