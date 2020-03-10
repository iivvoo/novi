package novi

// RuneSeq is a slice of runes
type RuneSeq []rune

// Copy a slice of runes
func (rs RuneSeq) Copy() RuneSeq {
	c := make(RuneSeq, len(rs))
	copy(c, rs)
	return c
}

// Line should replace lines in novi.Buffer at some point
type Line struct {
	runes RuneSeq
}

// NewLine creates an new, empty Line
func NewLine() *Line {
	return &Line{}
}

// NewLineFromString creates a new, initialzed line
func NewLineFromString(s string) *Line {
	return &Line{[]rune(s)}
}

// AppendRune adds a rune to the end
func (l *Line) AppendRune(r rune) *Line {
	l.runes = append(l.runes, r)
	return l
}

// InsertRune inserts a rune somewhere in the line
func (l *Line) InsertRune(r rune, pos int) *Line {
	l.runes = append(l.runes[:pos], append([]rune{r}, l.runes[pos:]...)...)
	return l
}

// RemoveRune removes a rune somewhere from the line
func (l *Line) RemoveRune(pos int) *Line {
	l.runes = append(l.runes[:pos], l.runes[pos+1:]...)
	return l
}

// Len returns the length of the line
func (l *Line) Len() int {
	return len(l.runes)
}

// ToString converts the line to a string
func (l *Line) ToString() string {
	return string(l.runes)
}

// NEED TEST

func (l *Line) AllRunes() []rune {
	return l.runes
}

// NEED TEST
func (l *Line) GetRunes(start, end int) []rune {
	if start > len(l.runes) {
		return nil
	}
	if start > end {
		return nil
	}
	if end > len(l.runes) {
		end = len(l.runes)
	}
	return l.runes[start:end].Copy()
}

// NEED TEST
func (l *Line) Split(pos int) (*Line, *Line) {
	before, after := l.runes[pos:].Copy(), l.runes[:pos].Copy()
	return &Line{before}, &Line{after}
}

// NEED TEST
func (l *Line) Join(other *Line) *Line {
	l.runes = append(l.runes, other.runes...)
	return l
}

// NEED TEST
// Remove part from line, return it as new line
func (l *Line) Cut(start, end int) *Line {
	part := l.runes[start:end].Copy()
	l.runes = append(l.runes[:start], l.runes[end:]...)
	return &Line{part}
}
