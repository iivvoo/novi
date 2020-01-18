package termui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/iivvoo/ovim/logger"
	"github.com/iivvoo/ovim/ovim"
)

var log = logger.GetLogger("termui")

type TermUI struct {
	// internal
	Screen         tcell.Screen
	Editor         *ovim.Editor
	ViewportX      int
	ViewportY      int
	EditAreaWidth  int
	EditAreaHeight int
	Width          int
	Height         int

	Status string
}

func NewTermUI(Editor *ovim.Editor) *TermUI {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.Show()

	w, h := s.Size()
	// adjust for statusbars and box, to be fixed XXX
	tui := &TermUI{s, Editor, 0, 0, w - 1, h - 3, w, h, ""}
	return tui
}

func (t *TermUI) GetDimension() (int, int) {
	return t.EditAreaWidth, t.EditAreaHeight
}

func (t *TermUI) SetSize(width, height int) {
	if width == 0 || height == 0 {
		log.Printf("Can't set a width or height of 0: %d %d", width, height)
		return
	}
	t.EditAreaWidth = width
	t.EditAreaHeight = height
	log.Printf("EditArea width, heigth set to %d, %d", width, height)
}

func (t *TermUI) Finish() {
	t.Screen.Fini()
}

func (t *TermUI) Loop(c chan ovim.Event) {
	go func() {
		defer ovim.RecoverFromPanic(func() {
			t.Finish()
		})
		for {
			ev := t.Screen.PollEvent()

			if ev != nil {
				log.Printf("%+v", ev)
			}
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if k := MapTCellKey(ev); k != nil {
					c <- k
				}
			case *tcell.EventResize:
			}
			// how to decide if we need update?
		}
	}()
}

func (t *TermUI) SetStatus(status string) {
	t.Status = status
}

func (t *TermUI) DrawBox() {
	for y := 0; y < t.EditAreaHeight; y++ {
		t.Screen.SetContent(t.EditAreaWidth, y, '|', nil, tcell.StyleDefault)
	}
	for x := 0; x < t.EditAreaWidth; x++ {
		t.Screen.SetContent(x, t.EditAreaHeight, '-', nil, tcell.StyleDefault)
	}
	t.Screen.SetContent(t.EditAreaWidth, t.EditAreaHeight, '+', nil, tcell.StyleDefault)
}

func (t *TermUI) DrawStatusbar() {
	x := 0
	for _, r := range t.Status { // XXX May overflow
		t.Screen.SetContent(x, t.Height-1, r, nil, tcell.StyleDefault)
		x++
	}
	for x < t.EditAreaWidth {
		t.Screen.SetContent(x, t.Height-1, ' ', nil, tcell.StyleDefault)
		x++

	}
}

/*
 Rendering on the screen always starts at (0,0), but characters taken from
 the editor are from a specific viewport
*/
func (t *TermUI) Render() {
	t.Width, t.Height = t.Screen.Size()

	primaryCursor := t.Editor.Cursors[0]
	if primaryCursor.Pos > t.ViewportX+t.EditAreaWidth-1 {
		t.ViewportX = primaryCursor.Pos - (t.EditAreaWidth - 1)
	}
	if primaryCursor.Pos < t.ViewportX {
		t.ViewportX = primaryCursor.Pos
	}

	if primaryCursor.Line > t.ViewportY+t.EditAreaHeight-1 {
		t.ViewportY = primaryCursor.Line - (t.EditAreaHeight - 1)
	}
	if primaryCursor.Line < t.ViewportY {
		t.ViewportY = primaryCursor.Line
	}

	t.DrawStatusbar()

	t.DrawBox()

	/*
	 * Print the text within the current viewports, padding lines with `fillRune`
	 * to clear any remainders. THe latter is relevant when scrolling, for example
	 */
	y := 0
	for _, line := range t.Editor.Buffer.GetLines(t.ViewportY, t.ViewportY+t.EditAreaHeight) {
		x := 0
		for _, rune := range line.GetRunes(t.ViewportX, t.ViewportX+t.EditAreaWidth) {
			t.Screen.SetContent(x, y, rune, nil, tcell.StyleDefault)
			x++
		}
		for x < t.EditAreaWidth {
			t.Screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
			x++
		}
		y++
	}
	for y < t.EditAreaHeight {
		for x := 0; x < t.EditAreaWidth; x++ {
			t.Screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
		y++
	}
	// To make the cursor blink, show/hide it?
	for _, cursor := range t.Editor.Cursors {
		if cursor.Line != -1 {
			t.Screen.ShowCursor(cursor.Pos-t.ViewportX, cursor.Line-t.ViewportY)
		}
		// else probably show at (0,0)
	}
	t.Screen.Sync()
}
