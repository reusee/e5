package e5

import (
	"io"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	TestWrapFunc(t, Info("foo"))

	info := Info("foo %s", "bar")(io.EOF)
	errString := info.Error()
	if !strings.Contains(errString, "EOF") {
		t.Fatal()
	}
	if !strings.Contains(errString, "foo bar") {
		t.Fatal()
	}
	if !is(info, io.EOF) {
		t.Fatal()
	}
}
