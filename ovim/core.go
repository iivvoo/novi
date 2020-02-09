package ovim

type Emulation interface {
	HandleEvent(InputID, Event) bool
	GetStatus(int) string
	SetChan(chan EmuEvent)
}

// Don't all UI's need Render()?
type InputUI interface {
	AskInput(string) InputSource
	CloseInput(InputSource)
	UpdateInput(InputSource, string, int)
}

type StatusUI interface {
	SetStatus(string)
	SetError(string)
}

type UI interface {
	Finish()
	Loop(chan Event)
	Render()
	GetDimension() (int, int)
}

// Core glues Editor, UI and Emulation together, passing messages along as necessary
type Core struct {
	Editor    *Editor
	UI        UI
	Input     InputUI
	Status    StatusUI
	Emulation Emulation
}

func NewCore(e *Editor, ui UI, input InputUI, status StatusUI, em Emulation) *Core {
	return &Core{Editor: e, UI: ui, Input: input, Status: status, Emulation: em}
}

func (c *Core) Loop() {
	// One handler can add to the other channel, make sure they don't block
	uiChan := make(chan Event, 10)
	emuChan := make(chan EmuEvent, 10)

	ui2emu := map[InputSource]InputID{0: 0}
	emu2ui := map[InputID]InputSource{0: 0}

	c.Emulation.SetChan(emuChan)
	c.UI.Render()
	c.UI.Loop(uiChan)
main:
	for {
		width, _ := c.UI.GetDimension()
		status := c.Emulation.GetStatus(width)
		c.Status.SetStatus(status)
		c.UI.Render()
		select {

		case ev := <-uiChan:
			// Filter event on what emulation subscribes to
			// invoke plugins/extensions in some order

			switch e := ev.(type) {
			case *KeyEvent, *CharacterEvent:
				log.Printf("Event %v", e)
				id, ok := ui2emu[e.GetSource()]
				if !ok {
					log.Printf("Got event from unmapped source: %d", e.GetSource())
				} else if !c.Emulation.HandleEvent(id, e) {
					break main
				}
			}
		case ev := <-emuChan:
			switch e := ev.(type) {
			// other events we can handle here: quit, save file, open file
			case *AskInputEvent:
				id := c.Input.AskInput(e.Prompt)
				log.Printf("Received AskInputEvent: %s -> %d", e.Prompt, id)
				ui2emu[id] = e.ID
				emu2ui[e.ID] = id
			case *CloseInputEvent:
				log.Printf("Core: CloseEvent %d", e.ID)
				source := emu2ui[e.ID]
				c.Input.CloseInput(source)
			case *UpdateInputEvent:
				source := emu2ui[e.ID]
				c.Input.UpdateInput(source, e.Text, e.Pos)
			case *SaveEvent:
				log.Printf("SaveEvent %s %v", e.Name, e.Force)
				if err := c.Editor.SaveFile(e.Name, e.Force); err != nil {
					c.Status.SetError("Could not save: " + err.Error())
				}
			case *QuitEvent:
				if c.Editor.Buffer.Modified && !e.Force {
					c.Status.SetError("Unsaved changes, please save first or use q!")
				} else {
					break main
				}
			case *ErrorEvent:
				c.Status.SetError(e.Message)
				log.Printf("ErrorEvent %s", e.Message)
			}
		}
	}
}
