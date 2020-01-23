package ovim

import (
	"bufio"
	"os"

	"github.com/iivvoo/ovim/logger"
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
	filename string
	Buffer   *Buffer
	Cursors  Cursors
}

func NewEditor() *Editor {
	e := &Editor{Buffer: NewEmptyBuffer()}
	e.Cursors = append(e.Cursors, e.Buffer.NewCursor(-1, 0))
	return e
}

func (e *Editor) GetFilename() string {
	return e.filename
}

func (e *Editor) LoadFile(name string) {
	// reset everything

	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	e.Buffer = BufferFromFile(file)
	e.filename = name
	e.Buffer.Modified = false
}

// SaveFile saves the buffer to a file
func (e *Editor) SaveFile() {
	// todo: take from buffer, optionally ask for name
	if e.filename == "" {
		log.Println("No filename set on buffer, can't save")
		return
	}
	if err := CopyFile(e.filename, e.filename+".bak"); err != nil {
		log.Printf("Failed to make backup copy for %s: %v", e.filename, err)
		return
	}
	f, err := os.Create(e.filename)
	if err != nil {
		log.Printf("Failed to open/create %s: %v", e.filename, err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, line := range e.Buffer.Lines {
		if _, err := w.WriteString(string(line) + "\n"); err != nil {
			log.Printf("Failed to Write %s: %v", e.filename, err)
			return
		}

	}
	if err := w.Flush(); err != nil {
		log.Printf("Failed to Flush %s: %v", e.filename, err)
	}
}

// SetCursor sets the first cursor at a specific position
func (e *Editor) SetCursor(row, col int) {
	e.Cursors[0].Line = row
	e.Cursors[0].Pos = col
	e.Cursors[0].Validate()
}

var CursorMap = map[KeyType]CursorDirection{
	KeyLeft:  CursorLeft,
	KeyRight: CursorRight,
	KeyUp:    CursorUp,
	KeyDown:  CursorDown,
	KeyHome:  CursorBegin,
	KeyEnd:   CursorEnd,
}
