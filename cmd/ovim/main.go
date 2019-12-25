package main

import (
	"fmt"
	"os"

	"gitlab.com/iivvoo/ovim/ovim"
	termui "gitlab.com/iivvoo/ovim/ui/term"
)

func Loop(e *ovim.Editor, t *termui.TermUI, c chan ovim.Event) {
	quit := make(chan bool)

	go func() {

	loop:
		for {
			ev := <-c
			switch ev := ev.(type) {
			case *ovim.KeyEvent:
				switch ev.Key {
				case ovim.KeyEscape:
					quit <- true
					break loop
				case ovim.KeyEnter:
					e.AddLine()
				case ovim.KeyLeft, ovim.KeyRight, ovim.KeyUp, ovim.KeyDown:
					e.MoveCursor(ev.Key)
				default:
					panic(ev)
				}
			case *ovim.CharacterEvent:
				e.PutRuneAtCursors(ev.Rune)
			default:
				panic(ev)
			}
			first := e.Cursors[0]
			row, col := first.Line, first.Pos
			lines := len(e.Lines)
			t.SetStatus(fmt.Sprintf("Edit: r %d c %d lines %d", row, col, lines))
			t.RenderTerm()
		}
	}()
	<-quit
}
func start() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ovim filename")
		os.Exit(1)
	}

	editor := ovim.NewEditor()
	ui := termui.NewTermUI(editor)
	defer ovim.RecoverFromPanic(func() {
		ui.Finish()
	})

	fileName := os.Args[1]

	editor.LoadFile(fileName)
	editor.SetCursor(8, 0)

	c := make(chan ovim.Event)

	ui.RenderTerm()
	ui.Loop(c)
	Loop(editor, ui, c)
	ui.Finish()
}

func main() {
	start()
}
