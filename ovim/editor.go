package ovim

import (
	"bufio"
	"errors"
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
	e := &Editor{Buffer: NewBuffer().InitializeEmptyBuffer()}
	e.Cursors = append(e.Cursors, e.Buffer.NewCursor(-1, 0))
	return e
}

func (e *Editor) GetFilename() string {
	return e.filename
}

func (e *Editor) LoadFile(name string) {
	// reset everything

	file, err := os.Open(name)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	e.Buffer.LoadFile(file)
	e.filename = name
	e.Buffer.Modified = false
}

var (
	ErrSaveNoName         = errors.New("No filename set")
	ErrSaveNoBackup       = errors.New("Could not create backup")
	ErrSaveFailedCreate   = errors.New("Could not create file")
	ErrSaveWrite          = errors.New("Failed to write contents")
	ErrSaveWouldOverwrite = errors.New("Not overwriting existing file")
	ErrSaveOther          = errors.New("Other error")
)

// SaveFile saves the buffer to a file
func (e *Editor) SaveFile(name string, force bool) error {
	changed := name != "" && e.filename != name
	exists := true

	if name != "" {
		e.filename = name
	}

	if _, err := os.Stat(e.filename); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			return ErrSaveOther
		}
	}

	// only overwrite existing file
	if changed && exists && !force {
		if _, err := os.Stat(e.filename); err == nil || !os.IsNotExist(err) {
			return ErrSaveWouldOverwrite
		}
	}

	if e.filename == "" {
		log.Println("No filename set on buffer, can't save")
		return ErrSaveNoName
	}

	if exists {
		if err := CopyFile(e.filename, e.filename+".bak"); err != nil {
			log.Printf("Failed to make backup copy for %s: %v", e.filename, err)
			return ErrSaveNoBackup
		}
	}
	f, err := os.Create(e.filename)
	if err != nil {
		log.Printf("Failed to open/create %s: %v", e.filename, err)
		return ErrSaveFailedCreate
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, line := range e.Buffer.Lines {
		if _, err := w.WriteString(string(line) + "\n"); err != nil {
			log.Printf("Failed to Write %s: %v", e.filename, err)
			return ErrSaveWrite
		}

	}
	if err := w.Flush(); err != nil {
		log.Printf("Failed to Flush %s: %v", e.filename, err)
	}

	e.Buffer.Modified = false
	return nil
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
