package termui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"gitlab.com/iivvoo/ovim/logger"
	"gitlab.com/iivvoo/ovim/ovim"
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

var KeyMap = map[tcell.Key]ovim.KeyType{
	// KeyBackspace2 is the 0x7F variant that's a regular character but > ' '
	tcell.KeyBackspace2: ovim.KeyBackspace,
	tcell.KeyEsc:        ovim.KeyEscape,
	tcell.KeyEnter:      ovim.KeyEnter,
	tcell.KeyUp:         ovim.KeyUp,
	tcell.KeyDown:       ovim.KeyDown,
	tcell.KeyLeft:       ovim.KeyLeft,
	tcell.KeyRight:      ovim.KeyRight,
	tcell.KeyHome:       ovim.KeyHome,
	tcell.KeyEnd:        ovim.KeyEnd,
	tcell.KeyPgUp:       ovim.KeyPgUp,
	tcell.KeyPgDn:       ovim.KeyPgDn,
	tcell.KeyBackspace:  ovim.KeyBackspace,
	tcell.KeyTab:        ovim.KeyTab,
	tcell.KeyDelete:     ovim.KeyDelete,
	tcell.KeyInsert:     ovim.KeyInsert,
	tcell.KeyF1:         ovim.KeyF1,
	tcell.KeyF2:         ovim.KeyF2,
	tcell.KeyF3:         ovim.KeyF3,
	tcell.KeyF4:         ovim.KeyF4,
	tcell.KeyF5:         ovim.KeyF5,
	tcell.KeyF6:         ovim.KeyF6,
	tcell.KeyF7:         ovim.KeyF7,
	tcell.KeyF8:         ovim.KeyF8,
	tcell.KeyF9:         ovim.KeyF9,
	tcell.KeyF10:        ovim.KeyF10,
	tcell.KeyF11:        ovim.KeyF11,
	tcell.KeyF12:        ovim.KeyF12,
}

type DecomposedKey struct {
	Modifier ovim.KeyModifier
	Key      ovim.KeyType
	Rune     rune
}

var DecomposeMap = map[tcell.Key]DecomposedKey{
	tcell.KeyCtrlSpace: DecomposedKey{ovim.ModCtrl, ovim.KeyRune, ' '},
	tcell.KeyCtrlA:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'a'},
	tcell.KeyCtrlB:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'b'},
	tcell.KeyCtrlC:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'c'},
	tcell.KeyCtrlD:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'd'},
	tcell.KeyCtrlE:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'e'},
	tcell.KeyCtrlF:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'f'},
	tcell.KeyCtrlG:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'g'},
	tcell.KeyCtrlH:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'h'},
	tcell.KeyCtrlI:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'i'},
	tcell.KeyCtrlJ:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'j'},
	tcell.KeyCtrlK:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'k'},
	tcell.KeyCtrlL:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'l'},
	tcell.KeyCtrlM:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'm'},
	tcell.KeyCtrlN:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'n'},
	tcell.KeyCtrlO:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'o'},
	tcell.KeyCtrlP:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'p'},
	tcell.KeyCtrlQ:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'q'},
	tcell.KeyCtrlR:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'r'},
	tcell.KeyCtrlS:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 's'},
	tcell.KeyCtrlT:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 't'},
	tcell.KeyCtrlU:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'u'},
	tcell.KeyCtrlV:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'v'},
	tcell.KeyCtrlW:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'w'},
	tcell.KeyCtrlX:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'x'},
	tcell.KeyCtrlY:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'y'},
	tcell.KeyCtrlZ:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'z'},
	// ctrl-Escape, according to tcell?
	tcell.KeyCtrlLeftSq:    DecomposedKey{ovim.ModCtrl, ovim.KeyEscape, ' '},
	tcell.KeyCtrlBackslash: DecomposedKey{ovim.ModCtrl, ovim.KeyRune, '\\'},
	// tcell.KeyCtrlRightSq:    DecomposedKey{ovim.ModCtrl, ovim.KeyRune, ' '},
	tcell.KeyCtrlCarat:      DecomposedKey{ovim.ModCtrl, ovim.KeyRune, '^'},
	tcell.KeyCtrlUnderscore: DecomposedKey{ovim.ModCtrl, ovim.KeyRune, '_'},
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
				/*
				   type EventKey struct {
				   	t   time.Time
				   	mod ModMask
				   	key Key
				   	ch  rune
				   }

				   rune = character, mod=alt, shift, etc
				   Key() is smashed, eg CtrlL

				   constants:
				   https://github.com/gdamore/tcell/blob/de7e78efa4a71b3f36c7154989c529dbdf9ae623/key.go#L83

				   use map to map types
				   Terminals don't (always) decompose modifiers. For example,
				   'H' is just H, not shift-h
				   but also a lot of contol characters are sent as-is,
				   so CtrlL is not ctrl-l

				   We can decompose this and transform CtrlL into ctrl-l

				*/
				key := ev.Key()
				if ovimKey, ok := KeyMap[key]; ok {
					c <- &ovim.KeyEvent{Key: ovimKey}
				} else if decomposed, ok := DecomposeMap[key]; ok {
					c <- &ovim.KeyEvent{Modifier: decomposed.Modifier, Key: decomposed.Key, Rune: decomposed.Rune}
				} else {
					c <- &ovim.CharacterEvent{Rune: ev.Rune()}
				}
			case *tcell.EventResize:
			}
			// how to decide if we need update?
		}
	}()
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

	t.DrawBox()

	/*
	 * Print the text within the current viewports, padding lines with `fillRune`
	 * to clear any remainders. THe latter is relevant when scrolling, for example
	 */
	fillRune := '~'
	for y, line := range t.Editor.Buffer.GetLines(t.ViewportY, t.ViewportY+t.EditAreaHeight) {
		x := 0
		for _, rune := range line.GetRunes(t.ViewportX, t.ViewportX+t.EditAreaWidth) {
			t.Screen.SetContent(x, y, rune, nil, tcell.StyleDefault)
			x++
		}
		for x < t.EditAreaWidth {
			t.Screen.SetContent(x, y, fillRune, nil, tcell.StyleDefault)
			x++
		}
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
