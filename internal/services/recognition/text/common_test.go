package text_test

import "errors"

type errorWriter struct{}

func (er *errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("something happened")
}
