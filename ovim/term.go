package ovim

import "github.com/gdamore/tcell"

type TermUI struct {
	Screen         tcell.Screen
	Editor         *Editor
	ViewportX      int
	ViewportY      int
	EditAreaWidth  int
	EditAreaHeight int
	Width          int
	Height         int
}

func NewTermUI(Screen tcell.Screen, Editor *Editor) *TermUI {
	return &TermUI{Screen, Editor, 0, 0, 40, 10, -1, -1}
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

/*
 Rendering on the screen always starts at (0,0), but characters taken from
 the editor are from a specific viewport
*/
func (t *TermUI) RenderTerm() {
	t.Width, t.Height = t.Screen.Size()
	// if smaller than estate we need, possibly adjust EditArea{Width.Height}

	// if any of the cursors is offscreen, adjust viewport
	// focus on primary cursor first
	primaryCursor := t.Editor.Cursors[0]
	if primaryCursor.Line > t.EditAreaHeight {
		t.ViewportY = (primaryCursor.Line - t.EditAreaHeight)
	}
	if primaryCursor.Pos > t.EditAreaWidth {
		t.ViewportX = (primaryCursor.Pos - t.EditAreaWidth)
	}

	y := 0

	t.DrawBox()

	// move slice magic, limit checks, to line
	for _, line := range t.Editor.Lines[t.ViewportY:] {
		x := 0
		startX := t.ViewportX
		endX := startX + t.EditAreaWidth
		if startX > len(line) {
			startX = 0
		}
		if endX > len(line) {
			endX = len(line)
		}
		for _, rune := range line[startX:endX] {
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
