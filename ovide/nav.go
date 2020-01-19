package ovide

import (
	"io/ioutil"
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

func NewNavTree(c chan Event) *NavTree {
	tree := &NavTree{c: c}
	tree.m = make(map[string]*NavTreeEntry)

	// get current path / folder name
	root := tview.NewTreeNode(".").
		SetColor(tcell.ColorRed)

	rootEntry := &NavTreeEntry{IsDir: true, FullPath: ".", Filename: ".", node: root}
	tree.TreeView = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.current = rootEntry

	tree.SetBorder(true)
	tree.SetTitle("Explorer")

	// Add the current directory to the root node.
	tree.add(root, ".")

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
			t.selectHandler(pNode)
		}
		// add fake entry
		tmp := tview.NewTreeNode("-> ____")
		pNode.AddChild(tmp)
		t.SetCurrentNode(tmp)
		t.temp = pNode

		t.c <- &NewFileEvent{ParentFolder: p}

	case 'f':
		// h for help?
	}
	return event
}

func (t *NavTree) RefreshTemp() {
	t.temp.ClearChildren()
	// a bit of a hack to abuse the handler here.
	t.selectHandler(t.temp)
	t.temp = nil
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

func (t *NavTree) add(target *tview.TreeNode, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fp := filepath.Join(path, file.Name())
		ref := &NavTreeEntry{IsDir: file.IsDir(), FullPath: fp, Filename: file.Name()}
		t.m[fp] = ref
		log.Printf("Setting map on entry path %s", fp)
		node := tview.NewTreeNode(file.Name()).SetReference(ref)
		ref.node = node
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
}
