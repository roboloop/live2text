package text_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/recognition/text"
)

func TestAddNonFinalText(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)

	formatter.Append("foo", false)
	require.Equal(t, "Foo\n", formatter.Format())

	formatter.Append("bar baz", false)
	require.Equal(t, "Bar baz\n", formatter.Format())

	formatter.Append("foo bar baz", false)
	require.Equal(t, "Foo bar\nbaz", formatter.Format())

	formatter.Append("foo bar baz abcdef", false)
	require.Equal(t, "Foo bar\nbaz abcdef", formatter.Format())

	formatter.Append("foo bar baz abcdefghi", false)
	require.Equal(t, "baz\nabcdefghi", formatter.Format())
}

func TestAddLongText(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)
	formatter.Append("foobarbazabcdefgeh", false)
	require.Equal(t, "Foobarbazabcdefgeh\n", formatter.Format())
}

func TestAddFinalText(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)

	formatter.Append("foo bar baz", true)
	require.Equal(t, "Foo bar\nbaz.", formatter.Format())

	formatter.Append("abc", true)
	require.Equal(t, "Foo bar\nbaz. Abc.", formatter.Format())

	formatter.Append("efg", true)
	require.Equal(t, "baz. Abc.\nEfg.", formatter.Format())

	formatter.Append("12345 67890 777", true)
	require.Equal(t, "Efg. 12345\n67890 777.", formatter.Format())
}

func TestEdgeWords(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)
	formatter.Append("123456789", true)
	formatter.Append("0", true)
	require.Equal(t, "123456789.\n0.", formatter.Format())

	formatter = text.NewSubtitleFormatter(2, 10)
	formatter.Append("123456789", true)
	formatter.Append("1 2 3 4 5", true)
	require.Equal(t, "123456789.\n1 2 3 4 5.", formatter.Format())

	formatter = text.NewSubtitleFormatter(2, 10)
	formatter.Append("123456789", true)
	formatter.Append("1 2 3 4 56", true)
	require.Equal(t, "1 2 3 4\n56.", formatter.Format())

	formatter = text.NewSubtitleFormatter(2, 10)
	formatter.Append("123456789", true)
	formatter.Append("1 2 3 4 5 6", false)
	require.Equal(t, "1 2 3 4 5\n6", formatter.Format())

	formatter.Append("1 2 3 4 5 6 7 8 9 0", false)
	require.Equal(t, "1 2 3 4 5\n6 7 8 9 0", formatter.Format())

	formatter.Append("1 2 3 4 5 6 7 8 9 0123", false)
	require.Equal(t, "6 7 8 9\n0123", formatter.Format())

	formatter.Append("1 2 3 4 5 6 7 8 9 012345", false)
	require.Equal(t, "6 7 8 9\n012345", formatter.Format())

	formatter.Append("1 2 3 4 5 6 7 8 9 01234", true)
	require.Equal(t, "6 7 8 9\n01234.", formatter.Format())

	formatter.Append("Abc", false)
	require.Equal(t, "6 7 8 9\n01234. Abc", formatter.Format())

	formatter.Append("Abc cde efg", false)
	require.Equal(t, "01234. Abc\ncde efg", formatter.Format())

	formatter.Append("abc cde efghijk", false)
	require.Equal(t, "cde\nefghijk", formatter.Format())

	formatter.Append("abc cde efghijk qwe asd zxc", false)
	require.Equal(t, "qwe asd\nzxc", formatter.Format())
}

func TestAddMixText(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)

	formatter.Append("foo", true)
	require.Equal(t, "Foo.\n", formatter.Format())

	formatter.Append("", true)
	require.Equal(t, "Foo.\n", formatter.Format())

	formatter.Append("bar baz abcedf ghijklmn", false)
	require.Equal(t, "baz abcedf\nghijklmn", formatter.Format())

	formatter.Append("foo", false)
	require.Equal(t, "Foo. Foo\n", formatter.Format())

	formatter.Append("This is final countdown", true)
	require.Equal(t, "is final\ncountdown.", formatter.Format())

	formatter.Append("No more", true)
	require.Equal(t, "countdown.\nNo more.", formatter.Format())
}

func TestNonAsciiText(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)

	formatter.Append("😊😊😊😊 😊😊😊😊", true)
	require.Equal(t, "😊😊😊😊 😊😊😊😊.\n", formatter.Format())

	formatter.Append("🤩 🤩 🤩 🤩", false)
	require.Equal(t, "😊😊😊😊 😊😊😊😊.\n🤩 🤩 🤩 🤩", formatter.Format())

	formatter.Append("🤩 🤩 🤩 🤩 🥳🥳🥳🥳", true)
	require.Equal(t, "🤩 🤩 🤩 🤩\n🥳🥳🥳🥳.", formatter.Format())
}

func TestNoSegments(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)
	require.Empty(t, formatter.Format())
}

func TestShowingLinesPersistent(t *testing.T) {
	t.Parallel()

	formatter := text.NewSubtitleFormatter(2, 10)

	formatter.Append("foo bar baz qux quux corge grault", false)
	require.Equal(t, "quux corge\ngrault", formatter.Format())

	formatter.Append("foo bar baz qux quux", false)
	require.Equal(t, "quux\n", formatter.Format())
}

func TestSubtitleWriter(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		formatter := text.NewSubtitleFormatter(0, 0)
		subtitleWriter := text.NewSubtitleWriter(formatter)

		err := subtitleWriter.PrintCandidate(0, "foo")
		require.NoError(t, err)
		require.Equal(t, "Foo\n", formatter.Format())

		err = subtitleWriter.PrintCandidate(0, "bar")
		require.NoError(t, err)
		require.Equal(t, "Bar\n", formatter.Format())

		err = subtitleWriter.PrintFinal(0, "baz")
		require.NoError(t, err)
		require.Equal(t, "Baz.\n", formatter.Format())

		err = subtitleWriter.Finalize()
		require.NoError(t, err)
		require.Equal(t, "Baz.\n", formatter.Format())
	})
}
