package viemu

import (
	"gitlab.com/iivvoo/ovim/ovim"
)

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
