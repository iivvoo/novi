package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"gitlab.com/iivvoo/ovim/ovim"
)

func main() {

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

	editor := ovim.NewEditor()

	text := []string{
		"Hello world",
		"",
		"So nice you see you",
		"this is a reeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeally long line",
		"bye now",
	}

	for _, textline := range text {
		editor.AddLine()

		for _, rune := range textline {
			editor.PutRuneAtCursors(rune)
		}
	}

	ui := ovim.NewTermUI(s, editor)

	// ui.RenderTerm()

	quit := make(chan struct{})

	/*
	 * overall structure:
	 * UI handles events (mouse, keys, etc) and sends generic events to the main loop,
	 * e.g. key-escape, enter, etc.
	 * Using mappings (and more) this is mapped to actions
	 */
	go func() {
		for {
			update := false
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyEnter:
					editor.AddLine()
					update = true
				case tcell.KeyCtrlL:
					update = true
				case tcell.KeyLeft:
					editor.MoveCursor(ovim.CursorLeft)
					update = true
				case tcell.KeyRight:
					editor.MoveCursor(ovim.CursorRight)
					update = true
				case tcell.KeyUp:
					editor.MoveCursor(ovim.CursorUp)
					update = true
				case tcell.KeyDown:
					editor.MoveCursor(ovim.CursorDown)
					update = true
				default:
					editor.PutRuneAtCursors(ev.Rune())
					update = true
				}
			case *tcell.EventResize:
				update = true
			}
			if update {
				ui.RenderTerm()
			}
		}
	}()

	<-quit

	s.Fini()
}
