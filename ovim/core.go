package ovim

type Emulation interface {
	HandleEvent(Event) bool
	GetStatus(int) string
}

type UI interface {
	Finish()
	Loop(chan Event)
	SetStatus(string)
	Render()
	GetDimension() (int, int)
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
		width, _ := c.UI.GetDimension()
		status := c.Emulation.GetStatus(width)
		c.UI.SetStatus(status)
		c.UI.Render()
		ev := <-eventChan
		// Filter event on what emulation subscribes to
		// invoke plugins/extensions in some order
		if !c.Emulation.HandleEvent(ev) {
			break
		}
	}
}
