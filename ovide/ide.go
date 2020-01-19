package ovide

import (
	"github.com/gdamore/tcell"
	viemu "github.com/iivvoo/ovim/emu/vi"
	"github.com/iivvoo/ovim/logger"
	"github.com/iivvoo/ovim/ovim"
	"github.com/rivo/tview"
)

var log = logger.GetLogger("ovide")

/*
 *
 * Stuff to do:
 * - display statusbar
 * - refactor Tab into proper Primitive (agnostic of contents)
 *   with proper id to identify tab, select existing, buttons.
 * - refactor editor into proper Primitive: area + status
 */

type Event interface{}

type QuitEvent struct{}

type OpenFileEvent struct {
	Filename string
	FullPath string
}

// Run just starts everything
func Run() {
	c := make(chan Event)

	app := tview.NewApplication()

	layout := tview.NewFlex()

	list := FileTree(c)
	tabs := NewTabs()
	layout.AddItem(list, 25, 0, true)
	layout.AddItem(tabs, 0, 1, false)

	// TODO: Include some sort of "debugging" Box
	go func() {
		for {
			log.Printf("Waiting for command")
			ev := <-c
			log.Printf("Got command %T %v", ev, ev)

			// We don't know where we were called from so make sure
			// we wrap our update
			// We can open the file etc before calling QueueUpdateDraw,
			// only schedule AddTab there..
			app.QueueUpdateDraw(func() {
				e := ev // local copy
				switch e := e.(type) {
				case *OpenFileEvent:
					log.Printf("Opening tab for %s", e.Filename)
					editor := ovim.NewEditor()
					editor.LoadFile(e.FullPath)
					editor.SetCursor(0, 0)

					emu := viemu.NewVi(editor)

					app.SetFocus(tabs.AddTab(e.FullPath, e.Filename, NewOviPrimitive(editor, emu, e.Filename)))
					log.Println("Done opening tab")
				case *QuitEvent:
					app.Stop()
				}
			})

		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlT {
			app.SetFocus(list)
		}
		return event
	})

	c <- &OpenFileEvent{FullPath: "sample.txt", Filename: "sample.txt"}
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
