package ovim

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

type TermUI struct {
	// internal
	Screen         tcell.Screen
	Editor         *Editor
	ViewportX      int
	ViewportY      int
	EditAreaWidth  int
	EditAreaHeight int
	Width          int
	Height         int

	Status1 string
	Status2 string
}

func NewTermUI(Editor *Editor) *TermUI {
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

func (t *TermUI) Loop() {
	quit := make(chan struct{})

	/*
	 * overall structure:
	 * UI handles events (mouse, keys, etc) and sends generic events to the main loop,
	 * e.g. key-escape, enter, etc.
	 * Using mappings (and more) this is mapped to actions
	 */

	// TODO: Move editor logic to editor - only handle UI specific events
	go func() {
		defer RecoverFromPanic(func() {
			close(quit)
			t.Finish()
		})
		for {
			update := false
			ev := t.Screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyEnter:
					t.Editor.AddLine()
					update = true
				case tcell.KeyCtrlL:
					update = true
				case tcell.KeyLeft:
					t.Editor.MoveCursor(CursorLeft)
					update = true
				case tcell.KeyRight:
					t.Editor.MoveCursor(CursorRight)
					update = true
				case tcell.KeyUp:
					t.Editor.MoveCursor(CursorUp)
					update = true
				case tcell.KeyDown:
					t.Editor.MoveCursor(CursorDown)
					update = true
				default:
					t.Editor.PutRuneAtCursors(ev.Rune())
					update = true
				}
			case *tcell.EventResize:
				update = true
			}
			first := t.Editor.Cursors[0]
			r, c := first.Line, first.Pos
			lines := len(t.Editor.Lines)
			t.SetStatus(fmt.Sprintf("Edit: r %d c %d lines %d", r, c, lines))
			if update {
				t.RenderTerm()
			}
		}
	}()
	<-quit
	t.Finish()
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
func (t *TermUI) RenderTerm() {
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

/*
  Test area. Test vscode behaviour

  aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
  cccccccccccccccccccccccccccccccccccccccccccccc
  dddddddddddddddddddddddddddddddddddddddddddddd
  eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
  ffffffffffffffffffffffff
  ggggggggggggggggggggggg
  hhhhhhhhhhhhhhhhhhhhhhhh
*/
