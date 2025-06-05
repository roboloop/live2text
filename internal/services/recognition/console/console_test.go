package console_test

import (
	"bytes"
	"testing"

	"live2text/internal/services/recognition/console"
)

func TestWriter(t *testing.T) {
	var buf bytes.Buffer

	var expected []byte
	writer := console.NewWriter(&buf)

	writer.PrintFail("foo")
	expected = []byte("\033[0;31m" + "\033[K" + "foo" + "\r")
	if !bytes.HasSuffix(buf.Bytes(), expected) {
		t.Fatalf("Got %v, expected %v", buf.Bytes(), expected)
	}

	writer.PrintFail("foo bar")
	expected = []byte("\033[0;31m" + "\033[K" + "foo bar" + "\r")
	if !bytes.HasSuffix(buf.Bytes(), expected) {
		t.Fatalf("Got %v, expected %v", buf.Bytes(), expected)
	}

	writer.PrintSuccess("foo bar")
	expected = []byte("\033[0;32m" + "\033[K" + "foo bar" + "\n")
	if !bytes.HasSuffix(buf.Bytes(), expected) {
		t.Fatalf("Got %v, expected %v", buf.Bytes(), expected)
	}

	writer.PrintFail("baz")
	expected = []byte("\033[0;31m" + "\033[K" + "baz" + "\r")
	if !bytes.HasSuffix(buf.Bytes(), expected) {
		t.Fatalf("Got %v, expected %v", buf.Bytes(), expected)
	}
}
