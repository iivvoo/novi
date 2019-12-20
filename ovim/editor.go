package ovim

/*
 Hoe modeleer je een een editor?
 Een tekst bestaat in principe uit regels. Regels kunnen ingevoegd, verwijderd, ingekort/verlengd worden.

 Een editor heeft N cursors

 Zonder regels kan je ook geen cursors hebben (of cursor is op positie -1)

 multi cursor behaviour (vscode):
 allemaal op specifieke positie
 als regel korter is, dan eind van regel
 maar als mogelijk, restore naar originele positie
 cursor up/down verliest x-positie, alle cursos krijgen dan zelfde positie. Dus / -> |
 Presentatie:
 Een line wrapt aan het eind, of wordt getruncate

*/

type Cursor struct {
	Line int
	Pos  int
}

type Line []rune
type Editor struct {
	Lines   []Line
	Cursors []*Cursor
}

func NewEditor() *Editor {
	e := &Editor{}
	e.Cursors = append(e.Cursors, &Cursor{-1, 0})
	return e
}

func (e *Editor) AddLine() {
	e.Lines = append(e.Lines, Line{})
	e.Cursors[0].Line++
	e.Cursors[0].Pos = 0
}

// https://github.com/golang/go/wiki/SliceTricks

func (e *Editor) PutRuneAtCursors(r rune) {
	for _, cursor := range e.Cursors {
		line := e.Lines[cursor.Line]
		line = append(line[:cursor.Pos], append(Line{r}, line[cursor.Pos:]...)...)
		e.Lines[cursor.Line] = line
		cursor.Pos++
	}
}

type CursorMovement int

const (
	CursorUp CursorMovement = iota
	CursorDown
	CursorLeft
	CursorRight
)

func (e *Editor) MoveCursor(movement CursorMovement) {
	for _, cursor := range e.Cursors {
		switch movement {
		case CursorUp:
			if cursor.Line > 0 {
				cursor.Line--
				if cursor.Pos > len(e.Lines[cursor.Line]) {
					cursor.Pos = len(e.Lines[cursor.Line])
				}
			}
		case CursorDown:
			// weirdness because empty last line that we want to position on
			if cursor.Line < len(e.Lines)-1 {
				cursor.Line++
				if cursor.Pos > len(e.Lines[cursor.Line]) {
					cursor.Pos = len(e.Lines[cursor.Line])
				}
			}
		case CursorLeft:
			if cursor.Pos > 0 {
				cursor.Pos--
			}
		case CursorRight:
			if cursor.Pos < len(e.Lines[cursor.Line]) {
				cursor.Pos++
			}
		}
	}
}
