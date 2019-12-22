package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"gitlab.com/iivvoo/ovim/ovim"
)

func recoverFromPanic(cleanup func()) {
	/*
			A normal debug.Stack() at this point looks like this:
			goroutine 23 [running]:
				runtime/debug.Stack(0xc00012a000, 0x584f40, 0xc0000d0120)
		       /opt/go1.13.1/src/runtime/debug/stack.go:24 +0x9d
				main.start.func1.1(0xc0000882a0, 0x5c3de0, 0xc00012a000)
		       /projects/ovim/cmd/ovim/main.go:65 +0x8d
				panic(0x584f40, 0xc0000d0120)
		       /opt/go1.13.1/src/runtime/panic.go:679 +0x1b2
				gitlab.com/iivvoo/ovim/ovim.(*TermUI).RenderTerm(0xc0000b8700)
		       /projects/ovim/ovim/term.go:90 +0x487
				main.start.func1(0xc0000882a0, 0x5c3de0, 0xc00012a000, 0xc00008c720, 0xc0000b8700)
		       /projects/ovim/cmd/ovim/main.go:110 +0x229
				created by main.start

			We don't care about the trace up until the panic and the line after it,
			so we need to do some fancy counting to skip that. Also, it's nice to
			have the actual error at the bottom
	*/
	if r := recover(); r != nil {
		cleanup()
		s := strings.Split(string(debug.Stack()), "\n")
		panicCheck := 0
		fmt.Println(s[0])
		for _, entry := range s {
			if panicCheck == 2 {
				fmt.Println(entry)
			} else if strings.HasPrefix(entry, "panic(") {
				panicCheck = 1
			} else if panicCheck == 1 {
				panicCheck = 2
			}
		}
		fmt.Printf("%s\n", r)
	}

}
func start() {

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
		defer recoverFromPanic(func() {
			close(quit)
			s.Fini()
		})
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
			first := editor.Cursors[0]
			r, c := first.Line, first.Pos
			lines := len(editor.Lines)
			ui.SetStatus(fmt.Sprintf("Edit: r %d c %d lines %d", r, c, lines))
			if update {
				ui.RenderTerm()
			}
		}
	}()

	<-quit

	s.Fini()
}

func main() {
	start()
}
