package novide

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/iivvoo/novi/novi"
	termui "github.com/iivvoo/novi/ui/term"
	"github.com/rivo/tview"
)

// We're always just one source since the IDE can provide a separate input source
const (
	MainSource    novi.InputSource = 0
	CommandSource novi.InputSource = 1
)

// implement the OviPrimitive
type Ovi struct {
	tview.Primitive
	Editor     *novi.Editor
	ViewportX  int
	ViewportY  int
	editArea   tview.Primitive
	statusArea *tview.TextView
	Source     novi.InputSource
	c          chan novi.Event
	InputPos   int
	statusMsg  string
	errorMsg   string
}

func NewOviPrimitive(e *novi.Editor) tview.Primitive {
	x := tview.NewFlex().SetDirection(tview.FlexRow)

	editArea := tview.NewBox()
	statusArea := tview.NewTextView()
	statusArea.SetDynamicColors(true)
	x.AddItem(editArea, 0, 1, true)
	x.AddItem(statusArea, 1, 1, false)

	o := &Ovi{
		ViewportX:  0,
		ViewportY:  0,
		Primitive:  x,
		Editor:     e,
		editArea:   editArea,
		statusArea: statusArea,
		Source:     MainSource,
		c:          nil,
		InputPos:   -1}

	editArea.SetDrawFunc(o.TviewRender)
	editArea.SetInputCapture(o.HandleInput)
	o.UpdateStatus("")
	return o
}

func (o *Ovi) SetChan(c chan novi.Event) {
	o.c = c
}

func (o *Ovi) TviewRender(screen tcell.Screen, xx, yy, width, height int) (int, int, int, int) {
	ui := termui.NewTCellUI(screen, xx, yy, width, height)
	ui.RenderTCell(o.Editor)
	if o.Source == CommandSource {
		x, y, _, _ := o.statusArea.GetInnerRect()
		screen.ShowCursor(x+o.InputPos, y)
	}
	return 0, 0, 0, 0
}

func (o *Ovi) GetDimension() (int, int) {
	_, _, w, h := o.statusArea.GetRect()
	return w, h
}

func (o *Ovi) UpdateStatus(status string) {
	if o.Source == MainSource && o.errorMsg == "" {
		o.statusArea.SetText("[-:-]" + status)
	}
	o.statusMsg = status
}

func (o *Ovi) UpdateError(message string) {
	o.errorMsg = message

	o.statusArea.SetText(fmt.Sprintf("[white:red]%s[-:-]", message))
}

func (o *Ovi) UpdateInput(input string, pos int) {
	if o.Source == CommandSource {
		o.statusArea.SetText(input)
		o.InputPos = pos
	}
}

func (o *Ovi) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	if e := termui.MapTCellKey(event); e != nil {
		e.SetSource(o.Source)
		o.c <- e
		return nil
	}
	return event
}
