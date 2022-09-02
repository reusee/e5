package e5

import (
	"io"
	"testing"
)

func TestDebug(t *testing.T) {
	err := Debug("foo")(io.EOF)
	if err.Error() != io.EOF.Error() {
		t.Fatal()
	}
	ErrorLevel = DebugLevel
	err = Debug("foo")(io.EOF)
	if err.Error() != "foo\n"+io.EOF.Error() {
		t.Fatal()
	}
	ErrorLevel = InfoLevel
}
