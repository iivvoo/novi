package ovim

import "github.com/gdamore/tcell"

type TermUI struct {
	Screen    tcell.Screen
	Editor    *Editor
	ViewportX int
	ViewportY int
	Width     int
	Height    int
}

func NewTermUI(Screen tcell.Screen, Editor *Editor) *TermUI {
	return &TermUI{Screen, Editor, 0, 0, -1, -1}
}

/*
 Rendering on the screen always starts at (0,0), but characters taken from
 the editor are from a specific viewport
*/
func (t *TermUI) RenderTerm() {
	t.Width, t.Height = t.Screen.Size()
	// if any of the cursors is offscreen, adjust viewport
	// focus on primary cursor first
	primaryCursor := t.Editor.Cursors[0]
	if primaryCursor.Line > t.Height {
		t.ViewportY = (primaryCursor.Line - t.Height)
	}
	if primaryCursor.Pos > t.Width {
		t.ViewportX = (primaryCursor.Pos - t.Width)
	}

	y := 0

	for _, line := range t.Editor.Lines[t.ViewportY:] {
		x := 0
		startX := t.ViewportX
		if startX > len(line) {
			startX = 0
		}
		for _, rune := range line[startX:] {
			t.Screen.SetContent(x, y, rune, nil, tcell.StyleDefault)
			x++
		}
		y++
	}
	// To make the cursor blink, show/hide it?
	for _, cursor := range t.Editor.Cursors {
		if cursor.Line != -1 {
			t.Screen.ShowCursor(int(cursor.Pos), int(cursor.Line))
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
