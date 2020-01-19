package ovide

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

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
