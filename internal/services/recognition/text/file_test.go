package text_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/services/recognition/text"
)

func TestFileWriter(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		fileWriter := text.NewFileWriter(&buf)

		err := fileWriter.PrintCandidate(0, "foo")
		require.NoError(t, err)
		require.Empty(t, buf.String())

		err = fileWriter.PrintCandidate(0, "bar")
		require.NoError(t, err)
		require.Empty(t, buf.String())

		err = fileWriter.PrintFinal(0, "baz")
		require.NoError(t, err)
		require.Equal(t, "baz\n", buf.String())

		err = fileWriter.PrintCandidate(0, "abc")
		require.NoError(t, err)
		require.Equal(t, "baz\n", buf.String())

		err = fileWriter.Finalize()
		require.NoError(t, err)
		require.Equal(t, "baz\nabc\n", buf.String())
	})

	t.Run("already finalized", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		fileWriter := text.NewFileWriter(&buf)

		err := fileWriter.PrintFinal(0, "foo")
		require.NoError(t, err)
		require.Equal(t, "foo\n", buf.String())

		err = fileWriter.Finalize()
		require.NoError(t, err)
	})

	t.Run("cannot print a final message", func(t *testing.T) {
		t.Parallel()

		fileWriter := text.NewFileWriter(&errorWriter{})

		err := fileWriter.PrintFinal(0, "foo")
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot print a final message")
	})
}
