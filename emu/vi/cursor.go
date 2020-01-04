package viemu

import (
	"github.com/iivvoo/ovim/ovim"
)

// Move moves the cursor the way Vi would do it
func Move(c *ovim.Cursor, movement ovim.CursorDirection) {
	switch movement {
	case ovim.CursorUp:
		if c.Line > 0 {
			c.Line--
			if c.Pos > len(c.Buffer.Lines[c.Line]) {
				c.Pos = len(c.Buffer.Lines[c.Line])
			}
		}
	case ovim.CursorDown:
		// weirdness because empty last line that we want to position on
		if c.Line < len(c.Buffer.Lines)-1 {
			c.Line++
			if c.Pos > len(c.Buffer.Lines[c.Line]) {
				c.Pos = len(c.Buffer.Lines[c.Line])
			}
		}
	case ovim.CursorLeft:
		if c.Pos > 0 {
			c.Pos--
		}
	case ovim.CursorRight:
		if c.Pos < len(c.Buffer.Lines[c.Line]) {
			c.Pos++
		}
	case ovim.CursorBegin:
		c.Pos = 0
	case ovim.CursorEnd:
		c.Pos = len(c.Buffer.Lines[c.Line]) - 1
		if c.Pos < 0 {
			c.Pos = 0
		}
	}
}

// MoveMany moves the cursor more than one position, if possible
func MoveMany(c *ovim.Cursor, movement ovim.CursorDirection, count int) {
	for i := 0; i < count; i++ {
		Move(c, movement)
	}
}

// MoveCursorRune moves cursor based on hjkl
func (em *Vi) MoveCursorRune(r rune, count int) bool {
	m := map[rune]ovim.CursorDirection{
		'h': ovim.CursorLeft,
		'j': ovim.CursorDown,
		'k': ovim.CursorUp,
		'l': ovim.CursorRight,
	}
	for _, c := range em.Editor.Cursors {
		for i := 0; i < count; i++ {
			Move(c, m[r])
		}
	}
	return true
}

// HandleMoveCursors moves the cursors based on the given event
func (em *Vi) HandleMoveCursors(ev ovim.Event) bool {
	for _, c := range em.Editor.Cursors {
		Move(c, ovim.CursorMap[ev.(ovim.KeyEvent).Key])
	}
	return true
}
