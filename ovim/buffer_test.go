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

		AssertBufferMatch(t, b, []string{"before"})
	})
	t.Run("Test empty buffer, after", func(t *testing.T) {
		b := BuildBuffer()
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, []string{"after"})
	})
	t.Run("Test single line buffer, before", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "before", true); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, []string{"before", "single line"})
	})
	t.Run("Test single line buffer, after", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}
		AssertBufferMatch(t, b, []string{"single line", "after"})
	})
}
