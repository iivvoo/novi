package ovide

import (
	"io/ioutil"
	"path/filepath"

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

type TreeEntry struct {
	IsDir    bool
	FullPath string
	Filename string
}

func FileTree(c chan Event) tview.Primitive {
	// get current path / folder name
	root := tview.NewTreeNode(".").
		SetColor(tcell.ColorRed)

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetBorder(true)
	tree.SetTitle("Explorer")

	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			ref := &TreeEntry{IsDir: file.IsDir(), FullPath: filepath.Join(path, file.Name()), Filename: file.Name()}
			node := tview.NewTreeNode(file.Name()).SetReference(ref)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, ".")

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			entry := reference.(*TreeEntry)

			if entry.IsDir {
				add(node, entry.FullPath)
			} else {
				log.Printf("Opening file %s", entry.Filename)
				c <- &OpenFileEvent{FullPath: entry.FullPath, Filename: entry.Filename}
				log.Printf("Command sent, should have been handled")
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return tree
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
