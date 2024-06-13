package e5

import (
	"errors"
	"fmt"
	"hash/maphash"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

// Stacktrace represents call stack frames
type Stacktrace struct {
	hashSum uint64
}

var framesInfo sync.Map // uint64 -> []Frame

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
	v, ok := framesInfo.Load(s.hashSum)
	if !ok {
		panic("bad key")
	}
	frames := v.([]Frame)

	var b strings.Builder
	for i, frame := range frames {
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

var pcsPool = sync.Pool{
	New: func() any {
		slice := make([]uintptr, 128)
		return &slice
	},
}

var hasherPool = sync.Pool{
	New: func() any {
		return new(maphash.Hash)
	},
}

// WrapStacktrace wraps current stacktrace
var WrapStacktrace = WrapFunc(func(prev error) error {
	if prev == nil {
		return nil
	}
	if stacktraceIncluded(prev) {
		return prev
	}

	pcs := *pcsPool.Get().(*[]uintptr)
	defer func() {
		pcs = pcs[:cap(pcs)]
		pcsPool.Put(&pcs)
	}()

	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	hasher := hasherPool.Get().(*maphash.Hash)
	defer func() {
		hasher.Reset()
		hasherPool.Put(hasher)
	}()
	for _, pc := range pcs {
		hasher.Write(
			unsafe.Slice(
				(*byte)(unsafe.Pointer(&pc)),
				unsafe.Sizeof(pc),
			),
		)
	}
	sum := hasher.Sum64()

	if _, ok := framesInfo.Load(sum); !ok {
		// construct frame infos
		frames := make([]Frame, 0, n)
		runtimeFrames := runtime.CallersFrames(pcs[:n])
		for {
			frame, more := runtimeFrames.Next()
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
			frames = append(frames, Frame{
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
		framesInfo.LoadOrStore(sum, frames)
	}

	stacktrace := &Stacktrace{
		hashSum: sum,
	}
	err := Join(stacktrace, prev)
	return err
})

func stacktraceIncluded(err error) bool {
	return errors.As(err, new(*Stacktrace))
}

var errStacktrace = errors.New("stacktrace")
