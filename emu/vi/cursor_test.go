package viemu

import (
	"testing"

	"github.com/iivvoo/ovim/ovim"
)

func TestCursorMove(t *testing.T) {
	emu := SetupVi(ModeCommand,
		"First line, next one is empty",
		"",
		"short",
		"",
		"Third last line",
		"Second last line",
		"Last line")
	b := emu.Editor.Buffer

	// regular movements where pos doesn't have to change/adjust
	t.Run("Move cursor UP", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		emu.Move(c, ovim.CursorUp)

		ovim.AssertCursor(t, c, 4, 5)
	})
	t.Run("Move cursor Down", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		emu.Move(c, ovim.CursorDown)

		ovim.AssertCursor(t, c, 6, 5)
	})
	t.Run("Move cursor Left", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		emu.Move(c, ovim.CursorLeft)

		ovim.AssertCursor(t, c, 5, 4)
	})
	t.Run("Move cursor Right", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		emu.Move(c, ovim.CursorRight)

		ovim.AssertCursor(t, c, 5, 6)
	})
	t.Run("Move cursor Begin", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		emu.Move(c, ovim.CursorBegin)

		ovim.AssertCursor(t, c, 5, 0)
	})
	t.Run("Move cursor End", func(t *testing.T) {
		c := b.NewCursor(5, 5)
		emu.Move(c, ovim.CursorEnd)

		ovim.AssertCursor(t, c, 5, 15)
	})

	// cases where the line isn't long enough to preserve the cursor's pos
	t.Run("Move cursor UP with adjust", func(t *testing.T) {
		c := b.NewCursor(2, 3)
		emu.Move(c, ovim.CursorUp)

		ovim.AssertCursor(t, c, 1, 0)
	})

	t.Run("Move cursor Down with adjust", func(t *testing.T) {
		c := b.NewCursor(2, 3)
		emu.Move(c, ovim.CursorDown)

		ovim.AssertCursor(t, c, 3, 0)
	})

	// corner cases
	t.Run("Move cursor UP at boundary", func(t *testing.T) {
		c := b.NewCursor(0, 3)
		emu.Move(c, ovim.CursorUp)

		ovim.AssertCursor(t, c, 0, 3)
	})

	t.Run("Move cursor Down at boundary", func(t *testing.T) {
		c := b.NewCursor(6, 3)
		emu.Move(c, ovim.CursorDown)

		ovim.AssertCursor(t, c, 6, 3)
	})

	// left/right should actualy go to prev/next line, if possible
}
