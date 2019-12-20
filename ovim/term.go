package ovim

import "github.com/gdamore/tcell"

func RenderTerm(e *Editor, s tcell.Screen) {
	y := 0
	for _, line := range e.Lines {
		x := 0
		for _, rune := range line {
			s.SetContent(x, y, rune, nil, tcell.StyleDefault)
			x++
		}
		y++
	}
	// To make the cursor blink, show/hide it?
	for _, cursor := range e.Cursors {
		if cursor.Line != -1 {
			s.ShowCursor(int(cursor.Pos), int(cursor.Line))
		}
		// else probably show at (0,0)
	}
	s.Sync()
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
