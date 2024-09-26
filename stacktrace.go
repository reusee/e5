package e5

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unique"
)

type _PCs [64]uintptr

// Stacktrace represents call stack frames
type Stacktrace struct {
	handle unique.Handle[_PCs]
}

var _ error = new(Stacktrace)

// Error implements error interface
func (s *Stacktrace) Error() string {
	pcs := s.handle.Value()
	i := 0
	for pcs[i] > 0 {
		i++
	}

	buf := new(strings.Builder)

	frames := runtime.CallersFrames(pcs[:i])
	firstLine := true
	for {
		frame, more := frames.Next()

		if strings.HasPrefix(frame.Function, "github.com/reusee/e5.") &&
			!strings.HasPrefix(frame.Function, "github.com/reusee/e5.Test") {
			// internal funcs
			if !more {
				break
			}
			continue
		}

		if firstLine {
			fmt.Fprintf(buf, "$ %v %v:%v", frame.Function, frame.File, frame.Line)
			firstLine = false
		} else {
			fmt.Fprintf(buf, "\n& %v %v:%v", frame.Function, frame.File, frame.Line)
		}

		if !more {
			break
		}
	}

	return buf.String()
}

func (s *Stacktrace) Is(err error) bool {
	// ignore content
	if _, ok := err.(*Stacktrace); ok {
		return true
	}
	return false
}

// WrapStacktrace wraps current stacktrace
var WrapStacktrace = WrapFunc(func(prev error) error {
	if prev == nil {
		return nil
	}
	if stacktraceIncluded(prev) {
		return prev
	}

	var pcs _PCs
	runtime.Callers(2, pcs[:])

	stacktrace := &Stacktrace{
		handle: unique.Make(pcs),
	}
	err := Join(stacktrace, prev)
	return err
})

func stacktraceIncluded(err error) bool {
	var p *Stacktrace
	return errors.Is(err, p)
}

var errStacktrace = errors.New("stacktrace")
