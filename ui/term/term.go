package termui

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/iivvoo/novi/logger"
	"github.com/iivvoo/novi/novi"
)

var log = logger.GetLogger("termui")

const (
	MainSource       novi.InputSource = 0
	MagicInputSource novi.InputSource = 31337
)

type TermUI struct {
	// internal
	Screen         tcell.Screen
	Editor         *novi.Editor
	ViewportX      int
	ViewportY      int
	EditAreaWidth  int
	EditAreaHeight int
	Width          int
	Height         int

	Status string
	Error  string

	// extra input support
	Source   novi.InputSource
	inputPos int
	prompt   string
	input    string
}

func NewTermUI(Editor *novi.Editor) *TermUI {
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
	log.Printf("Term size w h: %d %d", w, h)
	// adjust for statusbars and box, to be fixed XXX
	tui := &TermUI{s, Editor, 0, 0, w, h, w, h, "", "", MainSource, 0, "", ""}
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

func (t *TermUI) AskInput(prompt string) novi.InputSource {
	if t.Source == MagicInputSource {
		log.Println("term ui doesn't support more than one additional input!")
		return -1
	}
	t.Source = MagicInputSource
	t.prompt = prompt
	t.inputPos = len(t.prompt)
	// We only support one additional input so we can just make up some magic number
	return MagicInputSource
}

func (t *TermUI) CloseInput(source novi.InputSource) {
	t.input = ""
	t.Source = MainSource
}

func (t *TermUI) UpdateInput(source novi.InputSource, s string, pos int) {
	t.input = s
	t.inputPos = len(t.prompt) + pos
}

func (t *TermUI) Loop(c chan novi.Event) {
	go func() {
		defer novi.RecoverFromPanic(func() {
			t.Finish()
		})
		for {
			ev := t.Screen.PollEvent()

			// if ev != nil {
			// 	log.Printf("%+v", ev)
			// }
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if k := MapTCellKey(ev); k != nil {
					k.SetSource(t.Source)
					c <- k
				}
			case *tcell.EventResize:
			}
			// how to decide if we need update?
		}
	}()
}

func (t *TermUI) SetStatus(status string) {
	if t.Source == MainSource {
		t.Status = status
	}
}

func (t *TermUI) SetError(message string) {
	t.Error = message

	go func() {
		time.Sleep(time.Second * 3)
		t.Error = ""
		t.Render()
	}()
}
func (t *TermUI) drawBottomRow(s string, error bool) {
	// Use the full available width to draw the row, but make sure
	// status is truncated if too long
	x := 0

	style := tcell.StyleDefault

	if error {
		style = style.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	}

	rowY := t.EditAreaHeight - 1
	width := t.EditAreaWidth

	if len(s) > width {
		s = s[:width]
	}

	for _, r := range s { // XXX May overflow
		t.Screen.SetContent(x, rowY, r, nil, style)
		x++
	}
	for x < t.EditAreaWidth {
		t.Screen.SetContent(x, rowY, ' ', nil, tcell.StyleDefault)
		x++
	}

}
func (t *TermUI) DrawStatusbar() {
	if t.Error != "" {
		t.drawBottomRow(t.Error, true)
	} else {
		t.drawBottomRow(t.Status, false)
	}
}

func (t *TermUI) DrawInput() {
	t.drawBottomRow(t.prompt+t.input, false)
	t.Screen.ShowCursor(t.inputPos, t.Height-1)
}

func (t *TermUI) drawGutter(start, end int, guttersize int) {
	// drawGutter should decide size, return it

	// we need an upper limit - not all rows my be in use
	for y := 0; y < t.EditAreaHeight-1; y++ {
		l := ""
		lineno := y + start
		if lineno < end {
			l = strconv.Itoa(start + y + 1)
		}
		for len(l) < guttersize-1 {
			l = " " + l
		}
		for x, r := range l {
			t.Screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
}

/*
 Rendering on the screen always starts at (0,0), but characters taken from
 the editor are from a specific viewport
*/
func (t *TermUI) Render() {
	/*
	 t.EditAreaWidth is the size of the editor, t.EditAreaHeight the height.
	 Should this be purely the editing size, or also include statusbar, line gutter?
	 It will be set to terminal H/W if not set, so it's actually the entire area to use
	*/

	guttersize := 4 // 3 for numbers, 1 space)

	t.Width, t.Height = t.Screen.Size()
	editWidth, editHeight := t.EditAreaWidth-guttersize, t.EditAreaHeight-1

	primaryCursor := t.Editor.Cursors[0]
	if primaryCursor.Pos > t.ViewportX+editWidth-1 {
		t.ViewportX = primaryCursor.Pos - (editWidth - 1)
	}
	if primaryCursor.Pos < t.ViewportX {
		t.ViewportX = primaryCursor.Pos
	}

	if primaryCursor.Line > t.ViewportY+editHeight-1 {
		t.ViewportY = primaryCursor.Line - (editHeight - 1)
	}
	if primaryCursor.Line < t.ViewportY {
		t.ViewportY = primaryCursor.Line
	}

	/*
	 * Print the text within the current viewports, padding lines with `fillRune`
	 * to clear any remainders. THe latter is relevant when scrolling, for example
	 */
	y := 0
	for _, line := range t.Editor.Buffer.GetLines(t.ViewportY, t.ViewportY+editHeight) {
		x := guttersize
		for _, rune := range line.GetRunes(t.ViewportX, t.ViewportX+editWidth) {
			t.Screen.SetContent(x, y, rune, nil, tcell.StyleDefault)
			x++
		}
		for x < editWidth {
			t.Screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
			x++
		}
		y++
	}
	// We draw the gutter now because y contains the number of lines actually drawn (may be
	// at EOF), but eventually we'll need to do the math upfront and even let drawGutter calc is
	// own size
	t.drawGutter(t.ViewportY, t.ViewportY+y, guttersize)

	for y < editHeight {
		for x := guttersize; x < editWidth+guttersize; x++ {
			t.Screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
		y++
	}
	// To make the cursor blink, show/hide it?
	for _, cursor := range t.Editor.Cursors {
		if cursor.Line != -1 {
			t.Screen.ShowCursor(cursor.Pos-t.ViewportX+guttersize, cursor.Line-t.ViewportY)
		}
		// else probably show at (0,0)
	}

	// DrawInput may draw a cursor that has to override the main one
	// (tcell doesn't support multiple cursors?)
	// (but we may able to simulate those)
	if t.Source == MainSource {
		t.DrawStatusbar()
	} else {
		t.DrawInput()
	}

	t.Screen.Sync()
}
