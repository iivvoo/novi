package ovim

import (
	"testing"
)

func TestInsertLine(t *testing.T) {
	// empty buffer, before, after, top, bottom, middle
	t.Run("Test empty buffer, before", func(t *testing.T) {
		b := BuildBuffer()
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "before", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, "before")
	})
	t.Run("Test empty buffer, after", func(t *testing.T) {
		b := BuildBuffer()
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, "after")
	})
	t.Run("Test single line buffer, before", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "before", true); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, "before", "single line")
	})
	t.Run("Test single line buffer, after", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}
		AssertBufferMatch(t, b, "single line", "after")
	})
}

func TestRemoveCharacters(t *testing.T) {

	// before cursor (excluding cursor)
	t.Run("Test before, full deletion", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 14) // cursor on e
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "0123efghijklmnopqrstuvwyz")
	})
	t.Run("Test before, how many too large", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 5) // cursor on 5
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "56789abcdefghijklmnopqrstuvwyz")
	})
	t.Run("Test before, nothing before cursor", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 0) // cursor on 0
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "0123456789abcdefghijklmnopqrstuvwyz")
	})
	t.Run("Test before, empty line", func(t *testing.T) {
		b := BuildBuffer("")
		c := b.NewCursor(0, 0) // cursor on 0
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "")
	})
	// after cursor (including)
	t.Run("Test after, full deletion", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 14) // cursor on e
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "0123456789abcdopqrstuvwyz")
	})
	t.Run("Test after, how many too large", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 29) // cursor on t
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "0123456789abcdefghijklmnopqrs")
	})
	t.Run("Test before, nothing after cursor", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 34) // cursor on z
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "0123456789abcdefghijklmnopqrstuvwy")
	})
	t.Run("Test after, empty line", func(t *testing.T) {
		b := BuildBuffer("")
		c := b.NewCursor(0, 0) // cursor on 0
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "")
	})
}
