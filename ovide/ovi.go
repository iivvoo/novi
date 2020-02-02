package ovide

import (
	"github.com/gdamore/tcell"
	"github.com/iivvoo/ovim/ovim"
	termui "github.com/iivvoo/ovim/ui/term"
	"github.com/rivo/tview"
)

// implement the OviPrimitive
type Ovi struct {
	tview.Primitive
	Editor     *ovim.Editor
	Emulation  ovim.Emulation
	ViewportX  int
	ViewportY  int
	editArea   tview.Primitive
	statusArea *tview.TextView
}

func NewOviPrimitive(e *ovim.Editor, emu ovim.Emulation, name string) tview.Primitive {
	x := tview.NewFlex().SetDirection(tview.FlexRow)

	editArea := tview.NewBox()
	statusArea := tview.NewTextView()
	x.AddItem(editArea, 0, 1, true)
	x.AddItem(statusArea, 1, 1, false)

	o := &Ovi{
		ViewportX:  0,
		ViewportY:  0,
		Primitive:  x,
		Editor:     e,
		Emulation:  emu,
		editArea:   editArea,
		statusArea: statusArea}

	editArea.SetDrawFunc(o.TviewRender)
	editArea.SetInputCapture(o.HandleInput)
	o.UpdateStatus()
	return o
}

func (o *Ovi) TviewRender(screen tcell.Screen, xx, yy, width, height int) (int, int, int, int) {
	primaryCursor := o.Editor.Cursors[0]
	if primaryCursor.Pos > o.ViewportX+width-1 {
		o.ViewportX = primaryCursor.Pos - (width - 1)
	}
	if primaryCursor.Pos < o.ViewportX {
		o.ViewportX = primaryCursor.Pos
	}

	if primaryCursor.Line > o.ViewportY+height-1 {
		o.ViewportY = primaryCursor.Line - (height - 1)
	}
	if primaryCursor.Line < o.ViewportY {
		o.ViewportY = primaryCursor.Line
	}

	// Statusbar: get it from IDE or include in area?

	y := 0
	for _, line := range o.Editor.Buffer.GetLines(o.ViewportY, o.ViewportY+height) {
		x := 0
		for _, rune := range line.GetRunes(o.ViewportX, o.ViewportX+width) {
			screen.SetContent(xx+x, yy+y, rune, nil, tcell.StyleDefault)
			x++
		}
		for x < width {
			screen.SetContent(xx+x, yy+y, ' ', nil, tcell.StyleDefault)
			x++
		}
		y++
	}
	for y < height {
		for x := 0; x < width; x++ {
			screen.SetContent(xx+x, yy+y, ' ', nil, tcell.StyleDefault)
		}
		y++
	}
	// To make the cursor blink, show/hide it?
	for _, cursor := range o.Editor.Cursors {
		if cursor.Line != -1 {
			screen.ShowCursor(xx+cursor.Pos-o.ViewportX, yy+cursor.Line-o.ViewportY)
		}
		// else probably show at (0,0)
	}
	// Leave nothing for other components
	return 0, 0, 0, 0
}

func (o *Ovi) UpdateStatus() {
	_, _, w, _ := o.statusArea.GetRect()
	o.statusArea.SetText(o.Emulation.GetStatus(w))
}

func (o *Ovi) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	if e := termui.MapTCellKey(event); e != nil {
		if o.Emulation.HandleEvent(0, e) {
			o.UpdateStatus()
			return nil
		}
	}
	return event
}
