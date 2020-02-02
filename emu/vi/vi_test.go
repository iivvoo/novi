package viemu

import (
	"testing"

	"github.com/iivvoo/ovim/ovim"
)

func SetupVi(mode ViMode, lines ...string) *Vi {
	editor := ovim.NewEditor()
	editor.Buffer.LoadStrings(lines)
	emu := NewVi(editor)
	emu.Mode = mode

	return emu
}

func SetupViAndCursor(mode ViMode, line, pos int, lines ...string) (*Vi, *ovim.Cursor) {
	emu := SetupVi(mode, lines...)
	cursor := emu.Editor.Cursors[0]
	cursor.Line, cursor.Pos = line, pos
	return emu, cursor
}

func TestVi(t *testing.T) {
	t.Run("Cursor movement at end in Command mode", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeCommand, 0, 4, "hello")
		vi.MoveCursorRune('l', 1)

		ovim.AssertCursor(t, cursor, 0, 4) // expect no move
	})
	t.Run("Cursor movement at end in Edit mode", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 4, "hello")
		vi.MoveCursorRune('l', 1)

		ovim.AssertCursor(t, cursor, 0, 5) // Should move past end
	})
	t.Run("'a' append character at end of line", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 4, "hello")
		vi.HandleInsertionKeys(&ovim.CharacterEvent{Rune: 'a'})

		ovim.AssertCursor(t, cursor, 0, 5) // Should move past end
	})
	t.Run("Going to command mode when cursor past end will adjust cursor", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 5, "hello")
		vi.HandleToModeCommand(&ovim.KeyEvent{Key: ovim.KeyEscape})

		ovim.AssertCursor(t, cursor, 0, 4) // Should now be on last char
	})
	t.Run("Going to command mode on empty line", func(t *testing.T) {
		vi, cursor := SetupViAndCursor(ModeEdit, 0, 0, "")
		vi.HandleToModeCommand(&ovim.KeyEvent{Key: ovim.KeyEscape})

		ovim.AssertCursor(t, cursor, 0, 0) // Can't get any smaller
	})
}
