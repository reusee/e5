package e5

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestStacktrace(t *testing.T) {
	TestWrapFunc(t, WrapStacktrace)

	trace := WrapStacktrace(io.EOF)
	ok, err := regexp.MatchString(
		`\$ e5.stacktrace_test.go:[0-9]+ .*/e5/ e5.TestStacktrace\n&.*\n&.*\nEOF`,
		trace.Error(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("got %s", trace.Error())
	}
	if !is(trace, io.EOF) {
		t.Fatal()
	}
}

func TestDeepStacktrace(t *testing.T) {
	var foo func(int) error
	foo = func(i int) error {
		if i < 128 {
			return foo(i + 1)
		}
		return WrapStacktrace(io.EOF)
	}
	err := foo(1)
	if !errors.Is(err, io.EOF) {
		t.Fatal()
	}
}

func TestStacktraceIncluded(t *testing.T) {
	err := Error{
		WrapStacktrace(io.EOF),
	}
	if !stacktraceIncluded(err) {
		t.Fatal()
	}
	err = Join(
		io.EOF,
		WrapStacktrace(io.EOF),
	)
	if !stacktraceIncluded(err) {
		t.Fatal()
	}
}

func TestJoinStacktrace(t *testing.T) {
	err := Join(
		WrapStacktrace(io.EOF),
		WrapStacktrace(io.EOF),
	)
	str := err.Error()
	if len(strings.Split(str, "$")) != 2 {
		t.Fatalf("got %s", str)
	}
}
