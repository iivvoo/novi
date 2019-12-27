package ovim

type Cursor struct {
	Line int
	Pos  int
}

type CursorDirection uint8

const (
	CursorLeft CursorDirection = iota
	CursorRight
	CursorUp
	CursorDown
)

type Cursors []*Cursor

func (cs Cursors) Move(b *Buffer, movement CursorDirection) {
	for _, c := range cs {
		switch movement {
		case CursorUp:
			if c.Line > 0 {
				c.Line--
				if c.Pos > len(b.Lines[c.Line]) {
					c.Pos = len(b.Lines[c.Line])
				}
			}
		case CursorDown:
			// weirdness because empty last line that we want to position on
			if c.Line < len(b.Lines)-1 {
				c.Line++
				if c.Pos > len(b.Lines[c.Line]) {
					c.Pos = len(b.Lines[c.Line])
				}
			}
		case CursorLeft:
			if c.Pos > 0 {
				c.Pos--
			}
		case CursorRight:
			if c.Pos < len(b.Lines[c.Line]) {
				c.Pos++
			}
		}
	}
}

// MoveEnd, MoveStart - move all cursors to end/start of line
