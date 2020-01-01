package viemu

import (
	"gitlab.com/iivvoo/ovim/logger"
	"gitlab.com/iivvoo/ovim/ovim"
)

/*
 * Lots of stuff to do. Start with basic non-ex (?) commands, controls:
 * insert: iIoO
 * regular character insertion in edit mode
 * copy/paste (non/term/mouse: y, p etc)
 * commands: d10d, c5w, 10x, etc.
 *   escape in command mode -> cancel current
 */

var log = logger.GetLogger("viemu")

type ViMode int

const (
	ModeAny ViMode = iota
	ModeEdit
	ModeCommand
)

type Dispatch struct {
	Mode    ViMode
	Event   ovim.Event
	Events  []ovim.Event
	Handler func(ovim.Event)
}

func (d Dispatch) Do(event ovim.Event, mode ViMode) bool {
	if event.Equals(d.Event) && (d.Mode == ModeAny || d.Mode == mode) {
		d.Handler(event)
		log.Printf("Individual match %v %v", event, d)
		return true
	}
	for _, e := range d.Events {
		log.Printf("Events match %v %v", event, e)
		if event.Equals(e) && (d.Mode == ModeAny || d.Mode == mode) {
			log.Printf("Multi match %v %v", event, d)
			d.Handler(e)
			return true
		}
	}
	return false
}

type Vi struct {
	Editor *ovim.Editor
	Mode   ViMode

	dispatch []Dispatch
}

func NewVi(e *ovim.Editor) *Vi {
	em := &Vi{Editor: e, Mode: ModeCommand}
	dispatch := []Dispatch{
		Dispatch{Mode: ModeEdit, Event: ovim.KeyEvent{Key: ovim.KeyEscape}, Handler: em.HandleToModeCommand},
		Dispatch{Mode: ModeCommand, Event: ovim.CharacterEvent{Rune: 'i'}, Handler: em.HandleToModeEdit},
		Dispatch{Mode: ModeAny, Events: []ovim.Event{
			ovim.KeyEvent{Key: ovim.KeyLeft},
			ovim.KeyEvent{Key: ovim.KeyRight},
			ovim.KeyEvent{Key: ovim.KeyUp},
			ovim.KeyEvent{Key: ovim.KeyDown}}, Handler: em.HandleMoveCursors},
		Dispatch{Mode: ModeCommand, Events: []ovim.Event{
			ovim.CharacterEvent{Rune: 'h'},
			ovim.CharacterEvent{Rune: 'j'},
			ovim.CharacterEvent{Rune: 'k'},
			ovim.CharacterEvent{Rune: 'l'}}, Handler: em.HandleMoveHJKLCursors},
	}
	em.dispatch = dispatch
	return em
}

func (em *Vi) HandleToModeEdit(ovim.Event) {
	em.Mode = ModeEdit
}
func (em *Vi) HandleToModeCommand(ovim.Event) {
	em.Mode = ModeCommand
}

func (em *Vi) HandleMoveHJKLCursors(ev ovim.Event) {
	r := ev.(ovim.CharacterEvent).Rune

	m := map[rune]ovim.CursorDirection{
		'h': ovim.CursorLeft,
		'j': ovim.CursorDown,
		'k': ovim.CursorUp,
		'l': ovim.CursorRight,
	}
	em.Editor.Cursors.Move(em.Editor.Buffer, m[r])
}
func (em *Vi) HandleMoveCursors(ev ovim.Event) {
	em.Editor.MoveCursor(ev.(ovim.KeyEvent).Key)
}

func (em *Vi) HandleEvent(event ovim.Event) bool {
	for _, d := range em.dispatch {
		log.Printf("Match %v against %v\n", event, d)
		if d.Do(event, d.Mode) {
			log.Printf("  .. match!")
			return true
		}
	}
	return false
}

func (em *Vi) GetStatus(width int) string {
	if em.Mode == ModeEdit {
		return "--INSERT-- 10,20 (fake)"
	}
	return "10,20 (fake)"
}
