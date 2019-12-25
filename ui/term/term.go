package termui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"gitlab.com/iivvoo/ovim/ovim"
)

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

	Status1 string
	Status2 string
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

	tui := &TermUI{s, Editor, 0, 0, 40, 10, -1, -1,
		"Status 1 bla bla",
		"Status 2 bla bla"}
	return tui
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

			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					c <- &ovim.KeyEvent{ovim.KeyEscape}
				case tcell.KeyEnter:
					c <- &ovim.KeyEvent{ovim.KeyEnter}
				// case tcell.KeyCtrlL:
				case tcell.KeyLeft:
					c <- &ovim.KeyEvent{ovim.KeyLeft}
				case tcell.KeyRight:
					c <- &ovim.KeyEvent{ovim.KeyRight}
				case tcell.KeyUp:
					c <- &ovim.KeyEvent{ovim.KeyUp}
				case tcell.KeyDown:
					c <- &ovim.KeyEvent{ovim.KeyDown}
				default:
					c <- &ovim.CharacterEvent{ev.Rune()}
				}
			case *tcell.EventResize:
			}
			// how to decide if we need update?
		}
	}()
	// t.Finish()
}

func (t *TermUI) SetStatus(status string) {
	t.Status2 = status
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

func (t *TermUI) DrawStatusbar(bar string, pos int) {
	for i, r := range bar {
		t.Screen.SetContent(i, t.Height+pos-1, r, nil, tcell.StyleDefault)
	}
}

func (t *TermUI) DrawStatusbars() {
	t.DrawStatusbar(t.Status1, -1)
	t.DrawStatusbar(t.Status2, 0)
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

	t.Status1 = fmt.Sprintf("Term: vp %d %d size %d %d", t.ViewportX, t.ViewportY,
		t.Width, t.Height)
	t.DrawStatusbars()
	y := 0

	t.DrawBox()

	// move slice magic, limit checks, to line
	startY := t.ViewportY
	endY := startY + t.EditAreaHeight
	if endY > len(t.Editor.Lines) {
		endY = len(t.Editor.Lines)
	}
	for _, line := range t.Editor.Lines[startY:endY] {
		x := 0
		for _, rune := range line.GetRunes(t.ViewportX, t.ViewportX+t.EditAreaWidth) {
			t.Screen.SetContent(x, y, rune, nil, tcell.StyleDefault)
			x++
		}
		for x < t.EditAreaWidth {
			t.Screen.SetContent(x, y, '~', nil, tcell.StyleDefault)
			x++
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
