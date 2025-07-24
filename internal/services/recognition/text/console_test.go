package text_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/services/recognition/text"
)

func TestConsoleWriter(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		consoleWriter := text.NewConsoleWriter(&buf)

		err := consoleWriter.PrintCandidate(0, "foo")
		require.NoError(t, err)
		require.Equal(t, "\033[0;31m"+"foo"+"\033[K"+"\r", buf.String())

		err = consoleWriter.PrintCandidate(0, "bar")
		require.NoError(t, err)
		require.Equal(t,
			"\033[0;31m"+"foo"+"\033[K"+"\r"+
				"\033[0;31m"+"bar"+"\033[K"+"\r",
			buf.String(),
		)

		err = consoleWriter.PrintFinal(0, "baz")
		require.NoError(t, err)
		require.Equal(t,
			"\033[0;31m"+"foo"+"\033[K"+"\r"+
				"\033[0;31m"+"bar"+"\033[K"+"\r"+
				"\033[0;32m"+"baz"+"\033[K"+"\n",
			buf.String(),
		)

		err = consoleWriter.Finalize()
		require.NoError(t, err)
		require.Equal(t,
			"\033[0;31m"+"foo"+"\033[K"+"\r"+
				"\033[0;31m"+"bar"+"\033[K"+"\r"+
				"\033[0;32m"+"baz"+"\033[K"+"\n",
			buf.String(),
		)
	})

	t.Run("cannot print a final message", func(t *testing.T) {
		t.Parallel()

		consoleWriter := text.NewConsoleWriter(&errorWriter{})

		err := consoleWriter.PrintFinal(0, "foo")
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot print a final message")
	})

	t.Run("cannot print a candidate message", func(t *testing.T) {
		t.Parallel()

		consoleWriter := text.NewConsoleWriter(&errorWriter{})

		err := consoleWriter.PrintCandidate(0, "foo")
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot print a candidate message")
	})
}
