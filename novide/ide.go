package ovide

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell"
	viemu "github.com/iivvoo/novi/emu/vi"
	"github.com/iivvoo/novi/logger"
	"github.com/iivvoo/novi/novi"
	"github.com/rivo/tview"
)

var log = logger.GetLogger("ovide")

/*
 *
 * Stuff to do:
 * Open new files [ok]
 * Create + open new file
 * Save file
 * Navigate tabs
 * Implement generic input modal (for vi :command, search, stuff)
 * - Generic Error Modal (e.g. file create)
 */

// BuildEditor builds/sets up the editor and UI
func BuildEditor(app *tview.Application, fullpath string, editor *novi.Editor, ch chan IDEEvent) *Ovi {
	emu := viemu.NewVi(editor)
	prim := NewOviPrimitive(editor).(*Ovi)
	ui := NewWrapper(app, prim)

	c := novi.NewCore(editor, ui, emu)

	go func() {
		c.Loop()
		ch <- &CloseTabEvent{FullPath: fullpath}
	}()
	return prim
}

// Run just starts everything
func Run() {
	c := make(chan IDEEvent)

	app := tview.NewApplication()

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	cols := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(cols, 0, 1, true)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	nav := NewNavTree(c, cwd)
	tabs := NewTabs()
	cols.AddItem(nav, 25, 0, true)
	cols.AddItem(tabs, 0, 1, false)

	debug := tview.NewTextView()
	layout.AddItem(debug, 4, 0, false)
	pages := tview.NewPages().
		AddPage("ide", layout, true, true)

		// This loop handles IDE UI events (e.g. tree, tab).
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
					editor := novi.NewEditor()
					editor.LoadFile(e.FullPath)
					editor.SetCursor(0, 0)

					// pass a more generic tab id in stead of full path?
					prim := BuildEditor(app, e.FullPath, editor, c)

					app.SetFocus(tabs.AddTab(e.FullPath, e.Filename, prim))
					log.Println("Done opening tab")
				case *NewFileEvent:
					NewFileModal(app, pages).
						SetSuccess(func(s string) {
							p := filepath.Join(e.ParentFolder, s)
							os.Create(p)
							// DUP!
							editor := novi.NewEditor()
							editor.LoadFile(p)
							editor.SetCursor(0, 0)

							// pass a more generic tab id in stead of full path?
							prim := BuildEditor(app, p, editor, c)

							app.SetFocus(tabs.AddTab(p, s, prim))
							nav.ClearPlaceHolder()
							nav.SelectPath(p)
							// Select added file
						}).
						SetCancel(func(s string) {
							nav.ClearPlaceHolder()
						})
				case *CloseTabEvent:
					tabs.CloseTab(e.FullPath)
					log.Printf("Tab %s closed", e.FullPath)
					// XXX further cleanup?
					// Focus other tab (previously selected - keep track)
					// if no tabs, select Tree
					app.SetFocus(nav)
				case *DebugEvent:
					fmt.Fprintln(debug, e.Msg)
				case *QuitEvent:
					app.Stop()
				}
			})
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlT {
			app.SetFocus(nav)
		}
		return event
	})

	c <- &OpenFileEvent{FullPath: "sample.txt", Filename: "sample.txt"}
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
