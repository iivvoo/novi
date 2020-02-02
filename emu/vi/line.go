package viemu

// Line should replace lines in ovim.Buffer at some point
type Line struct {
	runes []rune
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

// later: Join
