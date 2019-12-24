package ovim

import (
	"bufio"
	"os"
)

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

type Editor struct {
	Lines   []Line
	Cursors []*Cursor
}

func NewEditor() *Editor {
	e := &Editor{}
	e.Cursors = append(e.Cursors, &Cursor{-1, 0})
	return e
}

func (e *Editor) LoadFile(name string) {
	// reset everything

	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		e.Lines = append(e.Lines, []rune(scanner.Text()))
	}
}

func (e *Editor) AddLine() {
	e.Lines = append(e.Lines, Line{})
	e.Cursors[0].Line++
	e.Cursors[0].Pos = 0
}

// SetCursor sets the first cursor at a specific position
func (e *Editor) SetCursor(row, col int) {
	e.Cursors[0].Line = row
	e.Cursors[0].Pos = col

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

func (e *Editor) MoveCursor(movement KeyType) {
	for _, cursor := range e.Cursors {
		switch movement {
		case KeyUp:
			if cursor.Line > 0 {
				cursor.Line--
				if cursor.Pos > len(e.Lines[cursor.Line]) {
					cursor.Pos = len(e.Lines[cursor.Line])
				}
			}
		case KeyDown:
			// weirdness because empty last line that we want to position on
			if cursor.Line < len(e.Lines)-1 {
				cursor.Line++
				if cursor.Pos > len(e.Lines[cursor.Line]) {
					cursor.Pos = len(e.Lines[cursor.Line])
				}
			}
		case KeyLeft:
			if cursor.Pos > 0 {
				cursor.Pos--
			}
		case KeyRight:
			if cursor.Pos < len(e.Lines[cursor.Line]) {
				cursor.Pos++
			}
		}
	}
}
