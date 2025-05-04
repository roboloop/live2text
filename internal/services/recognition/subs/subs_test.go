package subs

import (
	"strings"
	"testing"
)

func TestAddNonFinalText(t *testing.T) {
	writer := NewWriter(2, 10)

	writer.AddSection("foo", false)
	assertEqual(t, "foo", writer.Format())

	writer.AddSection("bar baz", false)
	assertEqual(t, "bar baz", writer.Format())

	writer.AddSection("foo bar baz", false)
	assertEqual(t, "foo bar\nbaz", writer.Format())

	writer.AddSection("foo bar baz abcdef", false)
	assertEqual(t, "foo bar\nbaz abcdef", writer.Format())

	writer.AddSection("foo bar baz abcdefghi", false)
	assertEqual(t, "baz\nabcdefghi", writer.Format())
}

func TestAddLongText(t *testing.T) {
	writer := NewWriter(2, 10)
	writer.AddSection("foobarbazabcdefgeh", false)
	assertEqual(t, "foobarbazabcdefgeh", writer.Format())
}

func TestAddFinalText(t *testing.T) {
	writer := NewWriter(2, 10)

	writer.AddSection("foo bar baz", true)
	assertEqual(t, "foo bar\nbaz.", writer.Format())

	writer.AddSection("abc", true)
	assertEqual(t, "foo bar\nbaz. abc.", writer.Format())

	writer.AddSection("efg", true)
	assertEqual(t, "baz. abc.\nefg.", writer.Format())

	writer.AddSection("12345 67890 777", true)
	assertEqual(t, "efg. 12345\n67890 777.", writer.Format())
}

func TestEdgeWords(t *testing.T) {
	writer := NewWriter(2, 10)
	writer.AddSection("123456789", true)
	writer.AddSection("0", true)
	assertEqual(t, "123456789.\n0.", writer.Format())

	writer = NewWriter(2, 10)
	writer.AddSection("123456789", true)
	writer.AddSection("1 2 3 4 5", true)
	assertEqual(t, "123456789.\n1 2 3 4 5.", writer.Format())

	writer = NewWriter(2, 10)
	writer.AddSection("123456789", true)
	writer.AddSection("1 2 3 4 56", true)
	assertEqual(t, "1 2 3 4\n56.", writer.Format())

	writer = NewWriter(2, 10)
	writer.AddSection("123456789", true)
	writer.AddSection("1 2 3 4 5 6", false)
	assertEqual(t, "1 2 3 4 5\n6", writer.Format())

	writer.AddSection("1 2 3 4 5 6 7 8 9 0", false)
	assertEqual(t, "1 2 3 4 5\n6 7 8 9 0", writer.Format())

	writer.AddSection("1 2 3 4 5 6 7 8 9 0123", false)
	assertEqual(t, "6 7 8 9\n0123", writer.Format())

	writer.AddSection("1 2 3 4 5 6 7 8 9 012345", false)
	assertEqual(t, "6 7 8 9\n012345", writer.Format())

	writer.AddSection("1 2 3 4 5 6 7 8 9 01234", true)
	assertEqual(t, "6 7 8 9\n01234.", writer.Format())

	writer.AddSection("abc", false)
	assertEqual(t, "6 7 8 9\n01234. abc", writer.Format())

	writer.AddSection("abc cde efg", false)
	assertEqual(t, "01234. abc\ncde efg", writer.Format())

	writer.AddSection("abc cde efghijk", false)
	assertEqual(t, "cde\nefghijk", writer.Format())

	writer.AddSection("abc cde efghijk qwe asd zxc", false)
	assertEqual(t, "qwe asd\nzxc", writer.Format())
}

func TestAddMixText(t *testing.T) {
	writer := NewWriter(2, 10)

	writer.AddSection("foo", true)
	assertEqual(t, "foo.", writer.Format())

	writer.AddSection("bar baz abcedf ghijklmn", false)
	assertEqual(t, "baz abcedf\nghijklmn", writer.Format())

	writer.AddSection("foo", false)
	assertEqual(t, "foo. foo", writer.Format())

	writer.AddSection("This is final countdown", true)
	assertEqual(t, "is final\ncountdown.", writer.Format())

	writer.AddSection("No more", true)
	assertEqual(t, "countdown.\nNo more.", writer.Format())
}

// Helper function for assertion
func assertEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}
