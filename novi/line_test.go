package novi

import "testing"

func AssertLineEquals(t *testing.T, l *Line, expected string) {
	t.Helper()

	if l.ToString() != expected {
		t.Errorf("Expected '%s' but got '%s' in stead", expected, l.ToString())
	}
}
func TestLine(t *testing.T) {
	t.Run("Initialize empty", func(t *testing.T) {
		l := NewLine()
		AssertLineEquals(t, l, "")
	})
	t.Run("Initialize string", func(t *testing.T) {
		l := NewLineFromString("hello")
		AssertLineEquals(t, l, "hello")
	})
	t.Run("Append to empty string", func(t *testing.T) {
		l := NewLine()
		l.AppendRune('a')
		AssertLineEquals(t, l, "a")
	})
	t.Run("Insert at start", func(t *testing.T) {
		l := NewLineFromString("ello")
		l.InsertRune('H', 0)
		AssertLineEquals(t, l, "Hello")
	})
	t.Run("Insert at end", func(t *testing.T) {
		l := NewLineFromString("Hell")
		l.InsertRune('o', 4)
		AssertLineEquals(t, l, "Hello")
	})
	t.Run("Insert in between", func(t *testing.T) {
		l := NewLineFromString("Hllo")
		l.InsertRune('e', 1)
		AssertLineEquals(t, l, "Hello")
	})
	t.Run("Remove from start", func(t *testing.T) {
		l := NewLineFromString("!Hello")
		l.RemoveRune(0)
		AssertLineEquals(t, l, "Hello")
	})
	t.Run("Remove from end", func(t *testing.T) {
		l := NewLineFromString("Hello!")
		l.RemoveRune(5)
		AssertLineEquals(t, l, "Hello")
	})
	t.Run("Remove from in between", func(t *testing.T) {
		l := NewLineFromString("He!llo")
		l.RemoveRune(2)
		AssertLineEquals(t, l, "Hello")
	})
}
