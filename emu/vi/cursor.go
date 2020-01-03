package viemu

import (
	"github.com/iivvoo/ovim/ovim"
)

// Move moves the cursor the way Vi would do it
func Move(b *ovim.Buffer, c *ovim.Cursor, movement ovim.CursorDirection) {
	switch movement {
	case ovim.CursorUp:
		if c.Line > 0 {
			c.Line--
			if c.Pos > len(b.Lines[c.Line]) {
				c.Pos = len(b.Lines[c.Line])
			}
		}
	case ovim.CursorDown:
		// weirdness because empty last line that we want to position on
		if c.Line < len(b.Lines)-1 {
			c.Line++
			if c.Pos > len(b.Lines[c.Line]) {
				c.Pos = len(b.Lines[c.Line])
			}
		}
	case ovim.CursorLeft:
		if c.Pos > 0 {
			c.Pos--
		}
	case ovim.CursorRight:
		if c.Pos < len(b.Lines[c.Line]) {
			c.Pos++
		}
	case ovim.CursorBegin:
		c.Pos = 0
	case ovim.CursorEnd:
		c.Pos = len(b.Lines[c.Line]) - 1
		if c.Pos < 0 {
			c.Pos = 0
		}
	}
}

// HandleMoveHJKLCursors hjkl can be used as cursor keys in command mode
func (em *Vi) HandleMoveHJKLCursors(ev ovim.Event) {
	r := ev.(ovim.CharacterEvent).Rune

	m := map[rune]ovim.CursorDirection{
		'h': ovim.CursorLeft,
		'j': ovim.CursorDown,
		'k': ovim.CursorUp,
		'l': ovim.CursorRight,
	}
	for _, c := range em.Editor.Cursors {
		Move(em.Editor.Buffer, c, m[r])
	}
}

// HandleMoveCursors moves the cursors based on the given event
func (em *Vi) HandleMoveCursors(ev ovim.Event) {
	for _, c := range em.Editor.Cursors {
		Move(em.Editor.Buffer, c, ovim.CursorMap[ev.(ovim.KeyEvent).Key])
	}
}
