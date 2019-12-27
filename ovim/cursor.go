package ovim

/*
 * Cursors always (?) operate in the context of a buffer with lines,
 * and are only valid based on that buffer. Bind cursors to buffer
 * expliclty?
 */

// Cursor defines a position within a buffer
type Cursor struct {
	Line int
	Pos  int
}

// CursorDirection defines the direction a cursor can go
type CursorDirection uint8

const (
	CursorLeft CursorDirection = iota
	CursorRight
	CursorUp
	CursorDown
	CursorBegin
	CursorEnd
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
		case CursorBegin:
			c.Pos = 0
		case CursorEnd:
			c.Pos = len(b.Lines[c.Line])
		}
	}
}

// MoveEnd, MoveStart - move all cursors to end/start of line
