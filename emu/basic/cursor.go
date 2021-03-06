package basicemu

import (
	"github.com/iivvoo/novi/novi"
)

// Move moves the cursor in the specified direction
func Move(c *novi.Cursor, movement novi.CursorDirection) {
	switch movement {
	case novi.CursorUp:
		if c.Line > 0 {
			c.Line--
			if c.Pos > c.Buffer.Lines[c.Line].Len() {
				c.Pos = c.Buffer.Lines[c.Line].Len()
			}
		}
	case novi.CursorDown:
		// weirdness because empty last line that we want to position on
		if c.Line < len(c.Buffer.Lines)-1 {
			c.Line++
			if c.Pos > c.Buffer.Lines[c.Line].Len() {
				c.Pos = c.Buffer.Lines[c.Line].Len()
			}
		}
	case novi.CursorLeft:
		if c.Pos > 0 {
			c.Pos--
		} else if c.Line > 0 {
			c.Line--
			c.Pos = c.Buffer.Lines[c.Line].Len()
		}
	case novi.CursorRight:
		if c.Pos < c.Buffer.Lines[c.Line].Len() {
			c.Pos++
		} else if c.Line < len(c.Buffer.Lines)-1 {
			c.Line++
			c.Pos = 0
		}
	case novi.CursorBegin:
		c.Pos = 0
	case novi.CursorEnd:
		// move *past* the end
		c.Pos = c.Buffer.Lines[c.Line].Len()
	}
}
