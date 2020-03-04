package viemu

import (
	"github.com/iivvoo/novi/novi"
)

// Move moves the cursor the way Vi would do it
// if pastEnd is true we allow the cursor to be placed past the last character,
// which Vi does in edit mode
func (em *Vi) Move(c *novi.Cursor, movement novi.CursorDirection) {
	maxPos := func(l novi.Line) int {
		limit := len(c.Buffer.Lines[c.Line])
		if em.Mode == ModeCommand {
			return limit
		}
		return limit + 1
	}

	switch movement {
	case novi.CursorUp:
		if c.Line > 0 {
			c.Line--
		}
	case novi.CursorDown:
		// weirdness because empty last line that we want to position on
		if c.Line < len(c.Buffer.Lines)-1 {
			c.Line++
		}
	case novi.CursorLeft:
		if c.Pos > 0 {
			c.Pos--
		}
	case novi.CursorRight:
		if c.Pos < maxPos(c.Buffer.Lines[c.Line]) {
			c.Pos++
		}
	case novi.CursorBegin:
		c.Pos = 0
	case novi.CursorEnd:
		c.Pos = len(c.Buffer.Lines[c.Line]) - 1
		if c.Pos < 0 {
			c.Pos = 0
		}
	}
	if l := maxPos(c.Buffer.Lines[c.Line]); c.Pos >= l {
		c.Pos = l - 1
	}
	if c.Pos < 0 {
		c.Pos = 0
	}
}

// MoveMany moves the cursor more than one position, if possible
func (em *Vi) MoveMany(c *novi.Cursor, movement novi.CursorDirection, count int) {
	for i := 0; i < count; i++ {
		em.Move(c, movement)
	}
}

// MoveCursorRune moves cursor based on hjkl
func (em *Vi) MoveCursorRune(r rune, count int) bool {
	m := map[rune]novi.CursorDirection{
		'h': novi.CursorLeft,
		'j': novi.CursorDown,
		'k': novi.CursorUp,
		'l': novi.CursorRight,
	}
	for _, c := range em.Editor.Cursors {
		for i := 0; i < count; i++ {
			em.Move(c, m[r])
		}
	}
	return true
}

// HandleMoveCursors moves the cursors based on the given event
func (em *Vi) HandleMoveCursors(ev novi.Event) bool {
	for _, c := range em.Editor.Cursors {
		em.Move(c, novi.CursorMap[ev.(*novi.KeyEvent).Key])
	}
	return true
}
