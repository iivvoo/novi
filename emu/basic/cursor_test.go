package basicemu

import (
	"testing"

	"github.com/iivvoo/novi/novi"
)

func TestCursorMove(t *testing.T) {
	b := novi.BuildBuffer("First line, next one is empty", "", "short", "", "Third last line", "Second last line", "Last line")

	// regular movements where pos doesn't have to change/adjust
	t.Run("Move cursor UP", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		Move(c, novi.CursorUp)

		novi.AssertCursor(t, c, 4, 5)
	})
	t.Run("Move cursor Down", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		Move(c, novi.CursorDown)

		novi.AssertCursor(t, c, 6, 5)
	})
	t.Run("Move cursor Left", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		Move(c, novi.CursorLeft)

		novi.AssertCursor(t, c, 5, 4)
	})
	t.Run("Move cursor Right", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		Move(c, novi.CursorRight)

		novi.AssertCursor(t, c, 5, 6)
	})
	t.Run("Move cursor Begin", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		Move(c, novi.CursorBegin)

		novi.AssertCursor(t, c, 5, 0)
	})
	t.Run("Move cursor End", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		Move(c, novi.CursorEnd)

		novi.AssertCursor(t, c, 5, 16)
	})

	// cases where the line isn't long enough to preserve the cursor's pos
	t.Run("Move cursor UP with adjust", func(t *testing.T) {
		c := b.NewCursor(2, 3)
		Move(c, novi.CursorUp)

		novi.AssertCursor(t, c, 1, 0)
	})

	t.Run("Move cursor Down with adjust", func(t *testing.T) {
		c := b.NewCursor(2, 3)
		Move(c, novi.CursorDown)

		novi.AssertCursor(t, c, 3, 0)
	})

	// move to adjecent lines
	t.Run("Move cursor Left at begin", func(t *testing.T) {
		c := b.NewCursor(6, 0)
		Move(c, novi.CursorLeft)

		novi.AssertCursor(t, c, 5, 16)
	})
	t.Run("Move cursor Right at end", func(t *testing.T) {
		c := b.NewCursor(5, 16)
		Move(c, novi.CursorRight)

		novi.AssertCursor(t, c, 6, 0)
	})

	// corner cases
	t.Run("Move cursor UP at boundary", func(t *testing.T) {
		c := b.NewCursor(0, 3)
		Move(c, novi.CursorUp)

		novi.AssertCursor(t, c, 0, 3)
	})

	t.Run("Move cursor Down at boundary", func(t *testing.T) {
		c := b.NewCursor(6, 3)
		Move(c, novi.CursorDown)

		novi.AssertCursor(t, c, 6, 3)
	})

	t.Run("Move cursor Left at 0,0", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		Move(c, novi.CursorLeft)

		novi.AssertCursor(t, c, 0, 0)
	})

	t.Run("Move cursor Right at complete end", func(t *testing.T) {
		c := b.NewCursor(6, 9)
		Move(c, novi.CursorRight)

		novi.AssertCursor(t, c, 6, 9)
	})
}
