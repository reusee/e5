package e5

import (
	"io"
	"testing"
)

func TestError(t *testing.T) {
	err := Chain(
		Chain(
			io.EOF,
			Chain(
				io.ErrClosedPipe,
				io.ErrNoProgress,
			),
		),
		Chain(
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
