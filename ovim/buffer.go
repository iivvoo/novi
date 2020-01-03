package ovim

/*
 * Buffer contains the implementaion of the text being manipulated by the editor.
 * It consists of an array of variable length strings of runes that can be manpipulated
 * using one or more cursors which identify one or more positions within the buffer
 */

// https://github.com/golang/go/wiki/SliceTricks

// Line implements a sequence of Runes
type Line []rune

// GetRunes implements safe slicing with bounday checks
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
	return l[start:end]
}

type Buffer struct {
	Lines    []Line
	Modified bool
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Length() int {
	return len(b.Lines)
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
	}
	b.Modified = true
}

/* SplitLines
 *
 * Split lines at position of cursors.
 * This is tricky since it will create extra lines, which may affect cursors below
 *
 * XXX make this a single cursor op since this makes it very hard to update cursors
 *
 * If we return some generic modification detail, we may be able to automatically update
 * cursors?
 */
func (b *Buffer) SplitLines(cs Cursors) {
	linesAdded := 0
	for _, c := range cs {
		line := b.Lines[c.Line+linesAdded]
		before, after := line[c.Pos:], line[:c.Pos]
		b.Lines = append(b.Lines[:c.Line+linesAdded],
			append([]Line{after, before}, b.Lines[c.Line+linesAdded+1:]...)...)
		linesAdded++
	}
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
