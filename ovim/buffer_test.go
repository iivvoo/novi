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
		AssertBufferModified(t, b, true)
	})
	t.Run("Test empty buffer, after", func(t *testing.T) {
		b := BuildBuffer()
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, "after")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test single line buffer, before", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "before", true); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, "before", "single line")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test single line buffer, after", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}
		AssertBufferMatch(t, b, "single line", "after")
		AssertBufferModified(t, b, true)
	})
}

func TestRemoveCharacters(t *testing.T) {

	// before cursor (excluding cursor)
	t.Run("Test before, full deletion", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 14) // cursor on e
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "0123efghijklmnopqrstuvwyz")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test before, how many too large", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 5) // cursor on 5
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "56789abcdefghijklmnopqrstuvwyz")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test before, nothing before cursor", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 0) // cursor on 0
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "0123456789abcdefghijklmnopqrstuvwyz")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test before, empty line", func(t *testing.T) {
		b := BuildBuffer("")
		c := b.NewCursor(0, 0) // cursor on 0
		b.RemoveCharacters(c, true, 10)
		AssertBufferMatch(t, b, "")
		AssertBufferModified(t, b, true)
	})
	// after cursor (including)
	t.Run("Test after, full deletion", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 14) // cursor on e
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "0123456789abcdopqrstuvwyz")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test after, how many too large", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 29) // cursor on t
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "0123456789abcdefghijklmnopqrs")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test before, nothing after cursor", func(t *testing.T) {
		b := BuildBuffer("0123456789abcdefghijklmnopqrstuvwyz")
		c := b.NewCursor(0, 34) // cursor on z
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "0123456789abcdefghijklmnopqrstuvwy")
		AssertBufferModified(t, b, true)
	})
	t.Run("Test after, empty line", func(t *testing.T) {
		b := BuildBuffer("")
		c := b.NewCursor(0, 0) // cursor on 0
		b.RemoveCharacters(c, false, 10)
		AssertBufferMatch(t, b, "")
		AssertBufferModified(t, b, true)
	})
}

func TestRemoveBetweenCursors(t *testing.T) {
	makeBuf := func() *Buffer {
		return BuildBuffer("Line 0 7890123456789a",
			"Line 1 7890123456789b",
			"Line 2 7890123456789c",
			"Line 3 7890123456789d")
	}

	t.Run("Single line test", func(t *testing.T) {
		b := makeBuf()
		res := b.RemoveBetweenCursors(b.NewCursor(1, 3), b.NewCursor(1, 12))
		AssertBufferMatch(t, b, "Line 0 7890123456789a",
			"Lin3456789b",
			"Line 2 7890123456789c",
			"Line 3 7890123456789d")
		AssertBufferMatch(t, res, "e 1 789012")
		AssertBufferModified(t, b, true)
	})

	t.Run("Dual line test", func(t *testing.T) {
		b := makeBuf()
		res := b.RemoveBetweenCursors(b.NewCursor(1, 10), b.NewCursor(2, 10))
		AssertBufferMatch(t, b,
			"Line 0 7890123456789a",
			"Line 1 789123456789c",
			"Line 3 7890123456789d")
		AssertBufferMatch(t, res, "0123456789b", "Line 2 7890")
		AssertBufferModified(t, b, true)
	})

	t.Run("Multi line test", func(t *testing.T) {
		b := makeBuf()
		res := b.RemoveBetweenCursors(b.NewCursor(0, 10), b.NewCursor(3, 10))
		AssertBufferMatch(t, b, "Line 0 789123456789d")
		AssertBufferMatch(t, res,
			"0123456789a",
			"Line 1 7890123456789b",
			"Line 2 7890123456789c",
			"Line 3 7890")
		AssertBufferModified(t, b, true)
	})

	t.Run("Remove full lines", func(t *testing.T) {
		b := makeBuf()
		res := b.RemoveBetweenCursors(b.NewCursor(0, 21), b.NewCursor(1, 20))
		AssertBufferMatch(t, b, "Line 0 7890123456789a",
			"Line 2 7890123456789c",
			"Line 3 7890123456789d")
		AssertBufferMatch(t, res, "", "Line 1 7890123456789b")
		AssertBufferModified(t, b, true)

	})
}
