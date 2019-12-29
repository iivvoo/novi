package ovim

import (
	"bufio"
	"fmt"
	"os"

	"gitlab.com/iivvoo/ovim/logger"
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

var log = logger.GetLogger("editor")

type Editor struct {
	Buffer  *Buffer
	Cursors Cursors
}

func NewEditor() *Editor {
	e := &Editor{Buffer: NewBuffer()}
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
		e.Buffer.AddLine(Line(scanner.Text()))
	}
}

// SaveFile saves the buffer to a file
func (e *Editor) SaveFile() {
	// todo: take from buffer, optionally ask for name
	fname := "/tmp/out1.txt"
	f, err := os.Create(fname)
	if err != nil {
		log.Printf("Failed to open/create %s: %v", fname, err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, line := range e.Buffer.Lines {
		if _, err := w.WriteString(string(line) + "\n"); err != nil {
			log.Printf("Failed to Write %s: %v", fname, err)
			return
		}

	}
	if err := w.Flush(); err != nil {
		log.Printf("Failed to Flush %s: %v", fname, err)
	}
}

// SetCursor sets the first cursor at a specific position
func (e *Editor) SetCursor(row, col int) {
	e.Cursors[0].Line = row
	e.Cursors[0].Pos = col

}

func (e *Editor) MoveCursor(movement KeyType) {
	m := map[KeyType]CursorDirection{
		KeyLeft:  CursorLeft,
		KeyRight: CursorRight,
		KeyUp:    CursorUp,
		KeyDown:  CursorDown,
	}

	if direction, ok := m[movement]; ok {
		e.Cursors.Move(e.Buffer, direction)

	} else {
		panic(fmt.Sprintf("Can't map key %v to a direction", movement))
	}
}
