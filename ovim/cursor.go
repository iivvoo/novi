package ovim

/*
 * Cursors always (?) operate in the context of a buffer with lines,
 * and are only valid based on that buffer. Bind cursors to buffer
 * expliclty?
 *
 * Generically speaking, you can assume that a cursor C at (l, p) will have effect
 * on lines before/after l and characters before/after pos. Sorts of effects:
 *
 * - lines have been inserted before/after (meaning cursors need update)
 * - lines have been removed before/after (meaning cursors need update)
 * - the line itself has been removed
 * - characters have been inserted before/after
 * - characters have been removed before/after
 * - the line has split meaning relatively complex cursor updates
 *
 * Depending on those modifications, cursus before/after that C and possibly C
 * itself needs updating
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
		c.Move(b, movement)
	}
}

func (c *Cursor) Move(b *Buffer, movement CursorDirection) {
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

// MoveEnd, MoveStart - move all cursors to end/start of line
