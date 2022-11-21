package e5

import (
	"io"
	"testing"
)

func TestError(t *testing.T) {
	err := Join(
		Join(
			io.EOF,
			Join(
				io.ErrClosedPipe,
				io.ErrNoProgress,
			),
		),
		Join(
			io.ErrShortBuffer,
			io.ErrUnexpectedEOF,
		),
	)
	if !is(err, io.EOF) {
		t.Fatal()
	}
	if !is(err, io.ErrClosedPipe) {
		t.Fatal()
	}
	if !is(err, io.ErrNoProgress) {
		t.Fatal()
	}
	if !is(err, io.ErrShortBuffer) {
		t.Fatal()
	}
	if !is(err, io.ErrUnexpectedEOF) {
		t.Fatal()
	}
}

func TestWith(t *testing.T) {
	TestWrapFunc(t, With(io.EOF))
}

func TestJoinSame(t *testing.T) {
	err := Join(
		Join(
			io.EOF,
			io.EOF,
		),
		Join(
			io.EOF,
			io.EOF,
		),
	)
	if str := err.Error(); str != "EOF" {
		t.Fatalf("got %s", str)
	}
}
