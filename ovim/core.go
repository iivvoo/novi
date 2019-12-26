package ovim

import "fmt"

type Emulation interface {
	HandleEvent(Event) bool
}

type UI interface {
	Finish()
	Loop(chan Event)
	SetStatus(string)
	Render()
}

type Core struct {
	Editor    *Editor
	UI        UI
	Emulation Emulation
}

func NewCore(e *Editor, ui UI, em Emulation) *Core {
	return &Core{Editor: e, UI: ui, Emulation: em}
}

func (c *Core) Loop() {
	eventChan := make(chan Event)

	c.UI.Render()
	c.UI.Loop(eventChan)
	for {
		ev := <-eventChan
		// Filter event on what emulation subscribes to
		// invoke plugins/extensions in some order
		if !c.Emulation.HandleEvent(ev) {
			break
		}

		first := c.Editor.Cursors[0]
		row, col := first.Line, first.Pos
		lines := c.Editor.Buffer.Length()
		c.UI.SetStatus(fmt.Sprintf("Edit: r %d c %d lines %d", row, col, lines))
		c.UI.Render()
	}

}
