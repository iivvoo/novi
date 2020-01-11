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
// A cursor has access to the buffer it's bound to, but should only use this for checking/reading,
// never modifying. Buffer modifications belonv on *Buffer
type Cursor struct {
	Line   int
	Pos    int
	Buffer *Buffer
}

// NewCursor creates a new cursor
func NewCursor(b *Buffer, line, pos int) *Cursor {
	return &Cursor{Line: line, Pos: pos, Buffer: b}
}

// Validate verifies if the cursor is valid for the current buffer
// and adjusts if necessary
func (c *Cursor) Validate() bool {
	valid := true

	if c.Line >= c.Buffer.Length() {
		c.Line = c.Buffer.Length() - 1
		valid = false
	}

	// XXX if we assert that the buffer always has at least 1 (empty)
	// line then c.Line can't be < 0
	if c.Line < 0 {
		c.Line = 0
		c.Pos = 0
		valid = false
	} else if c.Pos >= len(c.Buffer.Lines[c.Line]) {
		c.Pos = len(c.Buffer.Lines[c.Line]) - 1
		if c.Pos < 0 {
			c.Pos = 0
		}
		valid = false
	}
	return valid
}

// CursorDirection defines the direction a cursor can go
type CursorDirection uint8

// The different directions a Cursor can move to
const (
	CursorLeft CursorDirection = iota
	CursorRight
	CursorUp
	CursorDown
	CursorBegin
	CursorEnd
)

// Cursors is a collection of cursors
type Cursors []*Cursor

// After returns all cursors that are "after" c
func (cs Cursors) After(c *Cursor) Cursors {
	var result Cursors
	for _, cc := range cs {
		if cc.Line > c.Line || (cc.Line == c.Line && cc.Pos > c.Pos) {
			result = append(result, cc)
		}
	}
	return result
}

// Before returns all cursors that are "before" c
func (cs Cursors) Before(c *Cursor) Cursors {
	var result Cursors
	for _, cc := range cs {
		if cc.Line < c.Line || (cc.Line == c.Line && cc.Pos < c.Pos) {
			result = append(result, cc)
		}
	}
	return result
}
