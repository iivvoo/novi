package ovim_test

import (
	"testing"

	"gitlab.com/iivvoo/ovim/ovim"
)

func AssertBufferMatch(t *testing.T, b *ovim.Buffer, expected []string) {
	if a, b := b.Length(), len(expected); a != b {
		t.Errorf("Buffer size %d does not match expected size %d", a, b)
	}
	for i, e := range expected {
		if i >= b.Length() {
			break
		}
		if got := string(b.Lines[i]); e != got {
			t.Errorf("First mismatch at line %d\nexpected: %s\ngot     : %s", i, e, got)
		}
	}
}

func BuildBuffer(lines ...string) *ovim.Buffer {
	b := ovim.NewBuffer()
	for _, l := range lines {
		b.Lines = append(b.Lines, []rune(l))
	}

	return b
}

func TestInsertLine(t *testing.T) {
	// empty buffer, before, after, top, bottom, middle
	t.Run("Test empty buffer, before", func(t *testing.T) {
		b := BuildBuffer()
		c := &ovim.Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "before", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, []string{"before"})
	})
	t.Run("Test empty buffer, after", func(t *testing.T) {
		b := BuildBuffer()
		c := &ovim.Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, []string{"after"})
	})
	t.Run("Test single line buffer, before", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &ovim.Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "before", true); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}

		AssertBufferMatch(t, b, []string{"before", "single line"})
	})
	t.Run("Test single line buffer, after", func(t *testing.T) {
		b := BuildBuffer("single line")
		c := &ovim.Cursor{Line: 0, Pos: 0}
		if res := b.InsertLine(c, "after", false); !res {
			t.Error("Expected InsertLine to succeed but it didn't")
		}
		AssertBufferMatch(t, b, []string{"single line", "after"})
	})
}
