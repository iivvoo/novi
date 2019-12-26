package ovim

/*
 * Buffer contains the implementaion of the text being manipulated by the editor.
 * It consists of an array of variable length strings of runes that can be manpipulated
 * using one or more cursors which identify one or more positions within the buffer
 */

// for simplicity, cursors are currently defined here but should move to a
// separate file

// https://github.com/golang/go/wiki/SliceTricks
type Cursor struct {
	Line int
	Pos  int
}

type Cursors []*Cursor

func (cs Cursors) Move(b *Buffer, movement KeyType) {
	for _, c := range cs {
		switch movement {
		case KeyUp:
			if c.Line > 0 {
				c.Line--
				if c.Pos > len(b.Lines[c.Line]) {
					c.Pos = len(b.Lines[c.Line])
				}
			}
		case KeyDown:
			// weirdness because empty last line that we want to position on
			if c.Line < len(b.Lines)-1 {
				c.Line++
				if c.Pos > len(b.Lines[c.Line]) {
					c.Pos = len(b.Lines[c.Line])
				}
			}
		case KeyLeft:
			if c.Pos > 0 {
				c.Pos--
			}
		case KeyRight:
			if c.Pos < len(b.Lines[c.Line]) {
				c.Pos++
			}
		}
	}
}

// Line implements a sequence of Runes
type Line []rune

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
	Lines []Line
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Length() int {
	return len(b.Lines)
}

// AddLine adds a line to the bottom of the buffer
func (b *Buffer) AddLine(line Line) {
	b.Lines = append(b.Lines, line)
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
}
