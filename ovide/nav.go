package ovide

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
	c       chan Event
	current *NavTreeEntry
	temp    *tview.TreeNode
	m       map[string]*NavTreeEntry
}

func NewNavTree(c chan Event, path string) *NavTree {
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
		// Get the current selection
		// Get the parent folder
		current := t.current
		p := current.FullPath
		pNode := current.node
		log.Printf("About to create a new node on %s", p)
		if !current.IsDir {
			p = filepath.Dir(current.FullPath)
			pNode = t.m[p].node
			log.Printf("Not a folder, getting parent -> %s", p)
		} else {
			// This is supposed to open the folder
			t.selectHandler(pNode)
			pNode.SetExpanded(true)
		}
		// add fake entry
		tmp := tview.NewTreeNode("-> ____")
		pNode.AddChild(tmp)

		t.c <- &NewFileEvent{ParentFolder: p}
		t.SetCurrentNode(tmp)
		t.temp = pNode
		return nil

	case 'f':
		// h for help?
	}
	return event
}

func (t *NavTree) ClearPlaceHolder() {
	t.temp.ClearChildren()
	// a bit of a hack to abuse the handler here.
	t.selectHandler(t.temp)
	t.temp = nil
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
