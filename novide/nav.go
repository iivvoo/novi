package novide

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type NavTreeEntry struct {
	IsDir    bool
	FullPath string
	Filename string
	node     *tview.TreeNode
}

type NavTree struct {
	*tview.TreeView
	c          chan IDEEvent
	current    *NavTreeEntry
	workFolder *NavTreeEntry
	m          map[string]*NavTreeEntry
}

func NewNavTree(c chan IDEEvent, path string) *NavTree {
	current := filepath.Base(path)

	tree := &NavTree{c: c}
	tree.m = make(map[string]*NavTreeEntry)

	// get current path / folder name
	root := tview.NewTreeNode(current).SetColor(tcell.ColorRed)

	rootEntry := &NavTreeEntry{IsDir: true, FullPath: ".", Filename: ".", node: root}
	tree.m[path] = rootEntry

	tree.TreeView = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.current = rootEntry

	tree.SetBorder(true)
	tree.SetTitle("Explorer")

	// Add the current directory to the root node.
	tree.add(root, path)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(tree.selectHandler)
	tree.SetChangedFunc(tree.changeHandler)
	tree.SetInputCapture(tree.inputHandler)

	return tree
}

func (t *NavTree) changeHandler(node *tview.TreeNode) {
	if entry, ok := node.GetReference().(*NavTreeEntry); ok {
		t.c <- &DebugEvent{Msg: "Changed to " + entry.FullPath}
		t.current = entry
	}
}

func (t *NavTree) inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'n':
		// work will point to the folder that will hold the new file
		work := t.current
		log.Printf("About to create a new node on %s", work.FullPath)
		if !work.IsDir {
			p := filepath.Dir(work.FullPath)
			work = t.m[p] // XXX assert exists
			log.Printf("Not a folder, getting parent -> %s", p)
		} else {
			// This is supposed to open the folder
			t.add(work.node, work.FullPath)
			work.node.SetExpanded(true)
		}
		// add fake entry
		tmp := tview.NewTreeNode("-> ____")
		work.node.AddChild(tmp)
		t.SetCurrentNode(tmp)
		t.workFolder = work

		t.c <- &NewFileEvent{ParentFolder: work.FullPath}
		return nil

	case 'f':
		// h for help?
	}
	return event
}

func (t *NavTree) ClearPlaceHolder() {
	// Clear folder and re-add children
	t.workFolder.node.ClearChildren()
	t.add(t.workFolder.node, t.workFolder.FullPath)
	t.workFolder = nil
}

func (t *NavTree) SelectPath(p string) {
	if e, ok := t.m[p]; ok {
		log.Printf("Node for path %s found, selecting", p)
		t.SetCurrentNode(e.node)
	} else {
		log.Printf("Could not find node for path %s", p)
	}
}
func (t *NavTree) selectHandler(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		entry := reference.(*NavTreeEntry)
		t.c <- &DebugEvent{"Select " + entry.FullPath}

		if entry.IsDir {
			t.add(node, entry.FullPath)
		} else {
			log.Printf("Opening file %s", entry.Filename)
			t.c <- &OpenFileEvent{FullPath: entry.FullPath, Filename: entry.Filename}
			log.Printf("Command sent, should have been handled")
		}
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}

func (t *NavTree) createTreeEntry(basepath string, file os.FileInfo) *NavTreeEntry {
	fp := filepath.Join(basepath, file.Name())
	ref := &NavTreeEntry{IsDir: file.IsDir(), FullPath: fp, Filename: file.Name()}
	t.m[fp] = ref
	log.Printf("Setting map on entry path %s", fp)
	return ref
}

func (t *NavTree) add(target *tview.TreeNode, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		ref := t.createTreeEntry(path, file)
		node := tview.NewTreeNode(file.Name()).SetReference(ref)
		ref.node = node
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
}
