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
	EditAreaWidth  int
	EditAreaHeight int

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
	tui := &TermUI{s, Editor, w, h, "", "", MainSource, 0, "", ""}
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

/*
 Rendering on the screen always starts at (0,0), but characters taken from
 the editor are from a specific viewport
*/
func (t *TermUI) Render() {
	ui := NewTCellUI(t.Screen, 0, 0, t.EditAreaWidth, t.EditAreaHeight)
	ui.RenderTCell(t.Editor)

	if t.Source == MainSource {
		ui.RenderTCellStatusbar(t.Error, t.Status)
	} else {
		ui.RenderTCellInput(t.prompt+t.input, t.inputPos)
	}

	t.Screen.Sync()
}

/*
 * Only novi needs the statusbar - in novide it's separate controls. Embed in stead?
 * It does mean the -1 compensation won't make sense
 */

type TCellUI struct {
	baseX, baseY, width, height int
	screen                      tcell.Screen
}

func NewTCellUI(screen tcell.Screen, baseX, baseY, width, height int) *TCellUI {
	return &TCellUI{baseX, baseY, width, height, screen}
}

func (t *TCellUI) RenderTCellInput(s string, inputPos int) {
	t.RenderTCellBottomRow(s, false)
	t.screen.ShowCursor(t.baseX+inputPos, t.baseY+t.height-1)
}

func (t *TCellUI) RenderTCellBottomRow(s string, error bool) {
	// Use the full available width to draw the row, but make sure
	// status is truncated if too long
	x := 0

	style := tcell.StyleDefault

	if error {
		style = style.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	}

	rowY := t.height - 1

	if len(s) > t.width {
		s = s[:t.width]
	}

	for _, r := range s { // XXX May overflow
		t.screen.SetContent(t.baseX+x, t.baseY+rowY, r, nil, style)
		x++
	}
	for x < t.width {
		t.screen.SetContent(t.baseX+x, t.baseY+rowY, ' ', nil, tcell.StyleDefault)
		x++
	}
}

func (t *TCellUI) RenderTCellStatusbar(err, status string) {
	if err != "" {
		t.RenderTCellBottomRow(err, true)
	} else {
		t.RenderTCellBottomRow(status, false)
	}
}
func (t *TCellUI) RenderTCellGutter(start, end int, guttersize int) {
	// drawGutter should decide size, return it,
	// should perhaps check if numbering is enabled

	// we need an upper limit - not all rows my be in use
	for y := 0; y < t.height-1; y++ {
		l := ""
		lineno := y + start
		if lineno < end {
			l = strconv.Itoa(start + y + 1)
		}
		for len(l) < guttersize-1 {
			l = " " + l
		}
		for x, r := range l {
			t.screen.SetContent(t.baseX+x, t.baseY+y, r, nil, tcell.StyleDefault)
		}
	}
}

func (t *TCellUI) RenderTCell(editor *novi.Editor) {
	guttersize := 4 // 3 for numbers, 1 space)

	editWidth, editHeight := t.width-guttersize, t.height-1

	primaryCursor := editor.Cursors[0]

	ViewportX, ViewportY := 0, 0

	if primaryCursor.Pos > ViewportX+editWidth-1 {
		ViewportX = primaryCursor.Pos - (editWidth - 1)
	}
	if primaryCursor.Pos < ViewportX {
		ViewportX = primaryCursor.Pos
	}

	if primaryCursor.Line > ViewportY+editHeight-1 {
		ViewportY = primaryCursor.Line - (editHeight - 1)
	}
	if primaryCursor.Line < ViewportY {
		ViewportY = primaryCursor.Line
	}

	/*
	 * Print the text within the current viewports, padding lines with `fillRune`
	 * to clear any remainders. THe latter is relevant when scrolling, for example
	 */
	y := 0
	for _, line := range editor.Buffer.GetLines(ViewportY, ViewportY+editHeight) {
		x := 0
		for _, rune := range line.GetRunes(ViewportX, ViewportX+editWidth) {
			t.screen.SetContent(t.baseX+x+guttersize, t.baseY+y, rune, nil, tcell.StyleDefault)
			x++
		}
		for x < editWidth {
			t.screen.SetContent(t.baseX+x+guttersize, t.baseY+y, ' ', nil, tcell.StyleDefault)
			x++
		}
		y++
	}

	t.RenderTCellGutter(ViewportY, ViewportY+y, guttersize)

	for y < editHeight {
		for x := 0; x < editWidth; x++ {
			t.screen.SetContent(t.baseX+x+guttersize, t.baseY+y, ' ', nil, tcell.StyleDefault)
		}
		y++
	}
	// To make the cursor blink, show/hide it?
	for _, cursor := range editor.Cursors {
		if cursor.Line != -1 {
			t.screen.ShowCursor(t.baseX+cursor.Pos-ViewportX+guttersize, t.baseY+cursor.Line-ViewportY)
		}
		// else probably show at (0,0)
	}
}
