package e5

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Stacktrace represents call stack frames
type Stacktrace struct {
	Frames []Frame
}

// Frame represents a call frame
type Frame struct {
	File     string
	Dir      string
	Pkg      string
	Function string
	PkgPath  string
	Line     int
}

var _ error = new(Stacktrace)

// Error implements error interface
func (s *Stacktrace) Error() string {
	var b strings.Builder
	for i, frame := range s.Frames {
		if i == 0 {
			b.WriteString("$ ")
		} else {
			b.WriteString("\n& ")
		}
		b.WriteString(fmt.Sprintf(
			"%s:%s:%d %s %s",
			frame.Pkg,
			frame.File,
			frame.Line,
			frame.Dir,
			frame.Function,
		))
	}
	return b.String()
}

func (s *Stacktrace) Is(err error) bool {
	// ignore content
	if _, ok := err.(*Stacktrace); ok {
		return true
	}
	return false
}

var pcsPool = newPool(
	64,
	func() *[]uintptr {
		bs := make([]uintptr, 128)
		return &bs
	},
)

// WrapStacktrace wraps current stacktrace
var WrapStacktrace = WrapFunc(func(prev error) error {
	if prev == nil {
		return nil
	}
	if stacktraceIncluded(prev) {
		return prev
	}

	stacktrace := new(Stacktrace)
	v, put := pcsPool.Get()
	defer put()
	pcs := *v

	n := runtime.Callers(2, pcs)
	stacktrace.Frames = make([]Frame, 0, n)
	frames := runtime.CallersFrames(pcs[:n])
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
		dir, file := filepath.Split(frame.File)
		mod, fn := path.Split(frame.Function)
		if i := strings.Index(dir, mod); i > 0 {
			dir = dir[i:]
		}
		pkg := fn[:strings.IndexByte(fn, '.')]
		pkgPath := mod + pkg
		stacktrace.Frames = append(stacktrace.Frames, Frame{
			File:     file,
			Dir:      dir,
			Line:     frame.Line,
			Pkg:      pkg,
			Function: fn,
			PkgPath:  pkgPath,
		})
		if !more {
			break
		}
	}

	err := Join(stacktrace, prev)
	return err
})

func stacktraceIncluded(err error) bool {
	return errors.As(err, new(*Stacktrace))
}

var errStacktrace = errors.New("stacktrace")

// DropFrame returns a WrapFunc that drop Frames matching fn.
// If there is no existed stacktrace, a new one will be created
func DropFrame(fn func(Frame) bool) WrapFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		var stacktrace *Stacktrace
		if !errors.As(err, &stacktrace) {
			err = WrapStacktrace(err)
			errors.As(err, &stacktrace)
		}
		newFrames := stacktrace.Frames[:0]
		for _, frame := range stacktrace.Frames {
			if fn(frame) {
				continue
			}
			newFrames = append(newFrames, frame)
		}
		stacktrace.Frames = newFrames
		return err
	}
}

func WrapStacktraceWithoutPackageName(names ...string) WrapFunc {
	m := make(map[string]bool)
	for _, name := range names {
		m[name] = true
	}
	return WrapStacktrace.With(DropFrame(func(f Frame) bool {
		return m[f.Pkg]
	}))
}
