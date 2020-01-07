package viemu

import (
	"testing"

	"github.com/iivvoo/ovim/ovim"
)

func AssertLinePos(t *testing.T, expLine, expPos, actLine, actPos int) {
	t.Helper()
	if expLine != actLine || expPos != actPos {
		t.Errorf("Expected (l,p) to be (%d, %d) but got (%d, %d)", expLine, expPos, actLine, actPos)
	}
}

func TestJumpWordForward(t *testing.T) {
	b := ovim.BuildBuffer("This is the first line.", "", "  leading space",
		"trailing space   ", "last line")

	t.Run("Find first from start", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 0, 5, l, p)
	})
	t.Run("Find empty line", func(t *testing.T) {
		c := b.NewCursor(0, 19) // on 'i' in line
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 1, 0, l, p)
	})
	t.Run("Start from empty line", func(t *testing.T) {
		c := b.NewCursor(1, 0)
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 2, 2, l, p)
	})
	t.Run("Find from middle of word", func(t *testing.T) {
		c := b.NewCursor(0, 13) // on 'i' in first
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 0, 18, l, p)
	})
	t.Run("Find on next line", func(t *testing.T) {
		c := b.NewCursor(2, 11) // on 'a' in space
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 3, 0, l, p)
	})
	t.Run("Find end", func(t *testing.T) {
		c := b.NewCursor(4, 5) // on 'l' in line
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 4, 8, l, p)
	})
}

// Test 'B' behaviour
func TestJumpWordBackward(t *testing.T) {
	b := ovim.BuildBuffer("This is the first line.", "", "  leading space",
		"trailing space   ", "last line")

	t.Run("Find first from end", func(t *testing.T) {
		c := b.NewCursor(4, 8)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 4, 5, l, p)
	})
	t.Run("Find empty line", func(t *testing.T) {
		c := b.NewCursor(2, 2) // on 'l' in leading
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 1, 0, l, p)
	})
	t.Run("Find from middle of word", func(t *testing.T) {
		c := b.NewCursor(0, 13) // on 'i' in first
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 12, l, p)
	})
	t.Run("Find on previous line", func(t *testing.T) {
		c := b.NewCursor(4, 0)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 3, 9, l, p) // 's' in space
	})
	t.Run("From/to first word on line", func(t *testing.T) {
		c := b.NewCursor(3, 3)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 3, 0, l, p)
	})
	t.Run("From/to first word on line", func(t *testing.T) {
		b := ovim.BuildBuffer("First line", "    second line with spaces")
		c := b.NewCursor(1, 3) // at space character
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 6, l, p)
	})
	t.Run("Jump from empty", func(t *testing.T) {
		c := b.NewCursor(1, 0)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 18, l, p)
	})
	t.Run("Jump previous line word continues", func(t *testing.T) {
		b := ovim.BuildBuffer("foo bar", "4. bla 123")
		c := b.NewCursor(1, 0) // where we ended the previous test
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 4, l, p)
	})
}

func TestJump(t *testing.T) {
	b := ovim.BuildBuffer("This..isa line/;-with? separators", "", "  leading space",
		"https://github.com/some/repo.git?foo=a", "last line")

	t.Run("Find first from start", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 4, l, p)
	})
	t.Run("Find first from interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 4)
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 6, l, p)
	})
	t.Run("Find interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 10) // 'l' in line
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 14, l, p)
	})
	t.Run("Find first from interpunction, skipping whitespace", func(t *testing.T) {
		c := b.NewCursor(0, 21) // the ? after with
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 23, l, p) // expect space to be skipped
	})
	t.Run("Find empty line", func(t *testing.T) {
		c := b.NewCursor(0, 23) // the 's' in separators
		l, p := JumpForward(b, c)

		AssertLinePos(t, 1, 0, l, p) // expect space to be skipped
	})
	t.Run("Find from empty line", func(t *testing.T) {
		c := b.NewCursor(1, 0) // the 's' in separators
		l, p := JumpForward(b, c)

		AssertLinePos(t, 2, 2, l, p) // expect space to be skipped
	})
	t.Run("Jump URL", func(t *testing.T) {
		c := b.NewCursor(3, 0)
		l, p := JumpForward(b, c)

		AssertLinePos(t, 3, 5, l, p) // expect space to be skipped
	})
	t.Run("Jump URL 2", func(t *testing.T) {
		c := b.NewCursor(3, 5) // where we ended the previous test
		l, p := JumpForward(b, c)

		AssertLinePos(t, 3, 8, l, p)
	})
	t.Run("Jump next line word continues", func(t *testing.T) {
		b := ovim.BuildBuffer("foo bar", "bla 123")
		c := b.NewCursor(0, 4) // where we ended the previous test
		l, p := JumpForward(b, c)

		AssertLinePos(t, 1, 0, l, p)
	})
}
