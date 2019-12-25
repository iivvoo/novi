package ovim

import "fmt"

type Mapping interface {
}

type UI interface {
	Finish()
	Loop(chan Event)
	SetStatus(string)
	Render()
}

type Core struct {
	Editor  *Editor
	UI      UI
	Mapping Mapping
}

func NewCore(e *Editor, ui UI, m Mapping) *Core {
	return &Core{Editor: e, UI: ui, Mapping: m}
}

func (c *Core) Loop() {
	quit := make(chan bool)
	eventChan := make(chan Event)

	go func() {

	loop:
		for {
			ev := <-eventChan
			switch ev := ev.(type) {
			case *KeyEvent:
				switch ev.Key {
				case KeyEscape:
					quit <- true
					break loop
				case KeyEnter:
					c.Editor.AddLine()
				case KeyLeft, KeyRight, KeyUp, KeyDown:
					c.Editor.MoveCursor(ev.Key)
				default:
					panic(ev)
				}
			case *CharacterEvent:
				c.Editor.PutRuneAtCursors(ev.Rune)
			default:
				panic(ev)
			}
			first := c.Editor.Cursors[0]
			row, col := first.Line, first.Pos
			lines := len(c.Editor.Lines)
			c.UI.SetStatus(fmt.Sprintf("Edit: r %d c %d lines %d", row, col, lines))
			c.UI.Render()
		}
	}()

	c.UI.Render()
	c.UI.Loop(eventChan)
	<-quit
}
