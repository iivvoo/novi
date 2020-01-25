package ovim

type Emulation interface {
	HandleEvent(Event) bool
	GetStatus(int) string
	SetChan(chan Event)
}

type UI interface {
	Finish()
	Loop(chan Event)
	SetStatus(string)
	Render()
	GetDimension() (int, int)
	AskInput(string)
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
	// One handler can add to the other channel, make sure they don't block
	uiChan := make(chan Event, 2)
	emuChan := make(chan Event, 2)

	c.Emulation.SetChan(emuChan)
	c.UI.Render()
	c.UI.Loop(uiChan)
main:
	for {
		width, _ := c.UI.GetDimension()
		status := c.Emulation.GetStatus(width)
		c.UI.SetStatus(status)
		c.UI.Render()
		select {

		case ev := <-uiChan:
			// Filter event on what emulation subscribes to
			// invoke plugins/extensions in some order

			switch e := ev.(type) {
			case *KeyEvent, *CharacterEvent:
				if !c.Emulation.HandleEvent(e) {
					break main
				}
			}
		case ev := <-emuChan:
			switch e := ev.(type) {
			// other events we can handle here: quit, save file, open file
			case *AskInputEvent:
				log.Printf("Received AskInputEvent: %s", e.Prompt)
				c.UI.AskInput(e.Prompt)
			}
		}
	}
}
