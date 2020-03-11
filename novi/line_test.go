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
	t.Run("Get runes from line", func(t *testing.T) {
		l := NewLineFromString("123456789")
		r := l.GetRunes(3, 5)
		if len(r) != 2 {
			t.Errorf("Expected to get slice of size 2, got size %d", len(r))
		}

		if r[0] != '4' || r[1] != '5' {
			t.Errorf("Expected to get slice containing 45, got %v in stead", r)
		}
	})
	t.Run("Get runes from line, weird start/end", func(t *testing.T) {
		l := NewLineFromString("123456789")
		if r := l.GetRunes(100, 101); r != nil {
			t.Errorf("Expected to get nil, got %v", r)
		}

		if r := l.GetRunes(6, 4); r != nil {
			t.Errorf("Expected to get nil, got %v", r)
		}

		r := l.GetRunes(7, 20)
		if len(r) != 2 {
			t.Errorf("Expected to get slice of size 2, got size %d", len(r))
		}

		if r[0] != '8' || r[1] != '9' {
			t.Errorf("Expected to get slice containing 45, got %v in stead", r)
		}
	})

	t.Run("Get all runes", func(t *testing.T) {
		l := NewLineFromString("123456789")

		r := l.AllRunes()

		if len(r) != 9 {
			t.Errorf("Expected to get slice of size 9, got size %d", len(r))
		}

		if s := string(r); s != "123456789" {
			t.Errorf("Expected to get slice containing 1-9, got %s in stead", s)
		}
	})

	t.Run("Split line", func(t *testing.T) {
		l := NewLineFromString("123456789")

		l1, l2 := l.Split(5)

		AssertLineEquals(t, l1, "12345")
		AssertLineEquals(t, l2, "6789")
	})

	t.Run("Join line", func(t *testing.T) {
		l := NewLineFromString("1234")
		l2 := l.Join(NewLineFromString("56789"))
		AssertLineEquals(t, l, "123456789")
		AssertLineEquals(t, l2, "123456789")
	})

	t.Run("Cut line", func(t *testing.T) {
		l := NewLineFromString("123456789")
		c := l.Cut(3, 5)

		AssertLineEquals(t, c, "45")
		AssertLineEquals(t, l, "1236789")
	})
}
