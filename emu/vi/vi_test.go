package viemu

import (
	"testing"

	"github.com/iivvoo/novi/novi"
)

func SetupVi(mode ViMode, lines ...string) *Vi {
	editor := novi.NewEditor()
	editor.Buffer.LoadStrings(lines)
	emu := NewVi(editor)
	emu.Mode = mode

	return emu
}

func SetupViAndCursor(mode ViMode, line, pos int, lines ...string) (*Vi, *novi.Cursor) {
	emu := SetupVi(mode, lines...)
	cursor := emu.Editor.Cursors[0]
	cursor.Line, cursor.Pos = line, pos
	return emu, cursor
}

func TestVi(t *testing.T) {
	t.Run("Cursor movement at end in Command mode", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeCommand, 0, 4, "hello")
		vi.MoveCursorRune('l', 1)

		novi.AssertCursor(t, cursor, 0, 4) // expect no move
	})
	t.Run("Cursor movement at end in Edit mode", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 4, "hello")
		vi.MoveCursorRune('l', 1)

		novi.AssertCursor(t, cursor, 0, 5) // Should move past end
	})
	t.Run("'a' append character at end of line", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 4, "hello")
		vi.HandleInsertionKeys(&novi.CharacterEvent{Rune: 'a'})

		novi.AssertCursor(t, cursor, 0, 5) // Should move past end
	})
	t.Run("Going to command mode when cursor past end will adjust cursor", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 5, "hello")
		vi.HandleToModeCommand(&novi.KeyEvent{Key: novi.KeyEscape})

		novi.AssertCursor(t, cursor, 0, 4) // Should now be on last char
	})
	t.Run("Going to command mode on empty line", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 0, "")
		vi.HandleToModeCommand(&novi.KeyEvent{Key: novi.KeyEscape})

		novi.AssertCursor(t, cursor, 0, 0) // Can't get any smaller
	})
}
