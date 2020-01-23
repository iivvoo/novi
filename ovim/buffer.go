package ovim

import (
	"bufio"
	"io"
)

/*
 * Buffer contains the implementaion of the text being manipulated by the editor.
 * It consists of an array of variable length strings of runes that can be manpipulated
 * using one or more cursors which identify one or more positions within the buffer
 */

/*
  A buffer by default is empty, it has 0 lines. But cursors will point to (0,0) since they
  can't point anywhere else. Getting the line length at that point will crash

  we can also make sure a buffer always has one empty line
  perhaps add a flag? "start empty"?

  This way we can get rid of validate()?

  Either you start with an empty buffer, or you load data in the buffer. That's how you start.
  Empty buffer means start with single empty line, else load data UNLESS it's 0 bytes

  if you remove lines and you remove the last lines, replace it with an empty line

  Don't allow buffers to be created directly, always require it to be created and initialized
  - empty
  - from file
  - from lines

*/
// https://github.com/golang/go/wiki/SliceTricks

// Line implements a sequence of Runes
type Line []rune

// GetRunes implements safe slicing with boundary checks
func (l Line) GetRunes(start, end int) []rune {
	if start > len(l) {
		return nil
	}
	if start > end {
		return nil
	}
	if end > len(l) {
		end = len(l)
	}
	return l[start:end].Copy()
}

// Copy makes a proper copy of a line. Effectively, Lines are slices and taking
// (sub)slices of a slice does not make a copy. Modifying the original will copy the
// subslice which is not what you usually want
func (l Line) Copy() Line {
	c := make(Line, len(l))
	copy(c, l)
	return c
}

// Buffer encapsulates the state o an editable line buffer
type Buffer struct {
	Lines       []Line
	Modified    bool
	initialized bool
}

// NewBuffer creates a new Buffer. You usually don't want to call this directly
// since it will give you an unitialized buffer that you can't work with yet.
func NewBuffer() *Buffer {
	// the call you don't want since it doesn't initialize
	return &Buffer{}
}

func NewEmptyBuffer() *Buffer {
	return &Buffer{Lines: []Line{Line{}}, initialized: true}
}

func BufferFromFile(in io.Reader) *Buffer {
	b := NewBuffer()
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		b.AddLine(Line(scanner.Text()))
	}
	b.initialized = true
	return b
}

func BufferFromStrings(lines []string) *Buffer {
	b := NewBuffer()
	for _, l := range lines {
		b.Lines = append(b.Lines, []rune(l))
	}
	b.initialized = true
	return b
}

// NewCursor creates and binds a new cursor on this buffer
func (b *Buffer) NewCursor(line, pos int) *Cursor {
	return NewCursor(b, line, pos)
}

// Length returns the number of lines in this buffer
func (b *Buffer) Length() int {
	return len(b.Lines)
}

// Validate verifies and makes sure the buffer has a valid state
func (b *Buffer) Validate() bool {
	if b.Length() == 0 {
		b.Lines = []Line{Line{}}
		return false
	}
	return true
}

// GetLines attempts to retun the lines between start/end
func (b *Buffer) GetLines(start, end int) []Line {
	if start > b.Length() {
		return nil
	}
	if start > end {
		return nil
	}
	if end > b.Length() {
		end = b.Length()
	}
	return b.Lines[start:end]
}

// AddLine adds a line to the bottom of the buffer
func (b *Buffer) AddLine(line Line) {
	b.Lines = append(b.Lines, line)
	b.Modified = true
}

/* PutRuneAtCursor
 * Does not update cursors
 */
func (b *Buffer) PutRuneAtCursors(cs Cursors, r rune) {
	b.Validate()
	for _, c := range cs {
		line := b.Lines[c.Line]
		line = append(line[:c.Pos], append(Line{r}, line[c.Pos:]...)...)
		b.Lines[c.Line] = line
	}
	b.Modified = true
}

func (b *Buffer) RemoveRuneBeforeCursor(c *Cursor) {
	// We can't really do all cursors at once. Perhaps let caller always loop?
	// optionally, Cursors.all(func() {})
	if c.Pos > 0 {
		line := b.Lines[c.Line]
		line = append(line[:c.Pos-1], line[c.Pos:]...)
		b.Lines[c.Line] = line
		b.Modified = true
	}
}

/* SplitLine
 *
 * Split lines at position of cursors.
 * This is tricky since it will create extra lines, which may affect cursors below
 *
 * XXX make this a single cursor op since this makes it very hard to update cursors
 *
 * If we return some generic modification detail, we may be able to automatically update
 * cursors?
 */
func (b *Buffer) SplitLine(c *Cursor) {
	line := b.Lines[c.Line]
	before, after := line[c.Pos:], line[:c.Pos]
	b.Lines = append(b.Lines[:c.Line],
		append([]Line{after, before}, b.Lines[c.Line+1:]...)...)
	b.Modified = true
}

/* InsertLine
 *
 * insert line before/after given cursor
 */
func (b *Buffer) InsertLine(c *Cursor, line string, before bool) bool {
	// On an empty buffer, just add the line
	if b.Length() == 0 {
		b.AddLine([]rune(line))
		b.Modified = true
		return true
	}
	if c.Line >= b.Length() {
		return false
	}
	pos := c.Line
	if !before {
		pos++
	}
	b.Lines = append(b.Lines[:pos],
		append([]Line{[]rune(line)}, b.Lines[pos:]...)...)
	b.Modified = true
	return true
}

/* RemoveLine
 *
 * Remove and entire line. Does not update cursors wich may get invalid because
 * of the operation
 */
func (b *Buffer) RemoveLine(line int) bool {
	if line >= b.Length() {
		return false
	}
	b.Lines = append(b.Lines[:line], b.Lines[line+1:]...)
	b.Modified = true
	b.Validate()
	return true
}

/* JoinLineWithPrevious
 *
 * Join two lines: the one on the given position with the one before
 */
func (b *Buffer) JoinLineWithPrevious(line int) bool {
	if line == 0 || line > b.Length()-1 {
		return false
	}

	b.Lines[line-1] = append(b.Lines[line-1], b.Lines[line]...)
	b.RemoveLine(line)
	b.Modified = true
	return true
}

/* RemoveCharacters
 *
 * Removes a number of characters before or after the cursor
 * if after, it also removes the character the cursor is on.
 * if before, it preserves the character on the cursor. If this
 * is too vi-specific, refactor
 *
 * Could return number of characters actually removed?
 */
func (b *Buffer) RemoveCharacters(c *Cursor, before bool, howmany int) *Buffer {
	// RemoveCharacters implemented through RemoveBetweenCursors
	if before {
		startPos := c.Pos - howmany
		if startPos < 0 {
			startPos = 0
		}
		endPos := b.NewCursor(c.Line, c.Pos-1)
		return b.RemoveBetweenCursors(b.NewCursor(c.Line, startPos), endPos)
	}
	endPos := c.Pos + howmany - 1
	if endPos > len(b.Lines[c.Line])-1 {
		endPos = len(b.Lines[c.Line]) - 1
	}
	return b.RemoveBetweenCursors(c, b.NewCursor(c.Line, endPos))
}

func (b *Buffer) xRemoveCharacters(c *Cursor, before bool, howmany int) {
	line := b.Lines[c.Line]

	if before {
		if howmany > c.Pos {
			// befoe c.Pos, If c.Pos == 1, there's 1 char before it
			howmany = c.Pos
		}
		if howmany > 0 {
			line = append(line[:c.Pos-howmany], line[c.Pos:]...)
			b.Lines[c.Line] = line
		}
	} else {
		if howmany > len(line)-c.Pos {
			// 12345
			//   ^- cusor pos 2
			// max howmany: 3 -> len(line) - pos
			howmany = len(line) - c.Pos
		}
		if howmany > 0 {
			line = append(line[:c.Pos], line[c.Pos+howmany:]...)
			b.Lines[c.Line] = line
		}
	}
	b.Modified = true
}

// RemoveBetweenCursors removes all characters between start/end cursors (inclusive),
// across (entire) multiple lines if necessary. Returns the removed part as buffer
// Not suitable for block selections
func (b *Buffer) RemoveBetweenCursors(start, end *Cursor) *Buffer {
	res := &Buffer{}

	if start.Line > end.Line || (start.Line == end.Line && start.Pos > end.Pos) {
		return res
	}
	// We could check if start.IsBefore(end) and only act if true
	if end.Line > start.Line {
		first := b.Lines[start.Line][start.Pos:].Copy()
		res.Lines = append(res.Lines, first)

		middleSize := end.Line - start.Line - 1

		if middleSize > 0 {
			middle := b.Lines[start.Line+1 : end.Line]
			res.Lines = append(res.Lines, middle...)
		}

		last := b.Lines[end.Line][:end.Pos+1].Copy()
		res.Lines = append(res.Lines, last)

		b.Lines[start.Line] = append(b.Lines[start.Line][:start.Pos], b.Lines[end.Line][end.Pos+1:]...)

		// remove "end" line, since it was joined with start-line
		b.Lines = append(b.Lines[:end.Line], b.Lines[end.Line+1:]...)
		// now remove the middle part
		if middleSize > 0 {
			b.Lines = append(b.Lines[:start.Line+1], b.Lines[end.Line:]...)
		}
	} else { // removal is on same start/endline.
		part := b.Lines[start.Line][start.Pos : end.Pos+1].Copy()
		b.Lines[start.Line] = append(b.Lines[start.Line][:start.Pos], b.Lines[end.Line][end.Pos+1:]...)
		res.Lines = append(res.Lines, part)
	}
	b.Modified = true
	return res
}
