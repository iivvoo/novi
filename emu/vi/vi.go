package viemu

import (
	"fmt"

	"gitlab.com/iivvoo/ovim/logger"
	"gitlab.com/iivvoo/ovim/ovim"
)

/*
 * Lots of stuff to do. Start with basic non-ex (?) commands, controls:
 * insert: iIoOaA
 * backspace (similar behaviour as basic when joining lines)
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
		return true
	}
	for _, e := range d.Events {
		if event.Equals(e) && (d.Mode == ModeAny || d.Mode == mode) {
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
		Dispatch{Mode: ModeEdit, Event: ovim.CharacterEvent{}, Handler: em.HandleAnyRune},
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

// HandleAnyRune simply inserts the character in edit mode
func (em *Vi) HandleAnyRune(ev ovim.Event) {
	r := ev.(*ovim.CharacterEvent).Rune
	em.Editor.Buffer.PutRuneAtCursors(em.Editor.Cursors, r)
	em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorRight)
}

func (em *Vi) HandleMoveCursors(ev ovim.Event) {
	em.Editor.MoveCursor(ev.(ovim.KeyEvent).Key)
}

func (em *Vi) HandleEvent(event ovim.Event) bool {
	for _, d := range em.dispatch {
		log.Printf("Match %v against %v\n", event, d)
		if d.Do(event, em.Mode) {
			log.Printf("  .. match!")
			return true
		}
	}
	return false
}

func (em *Vi) GetStatus(width int) string {
	mode := ""
	first := em.Editor.Cursors[0]
	if em.Mode == ModeEdit {
		mode = "--INSERT-- "
	}
	return mode + fmt.Sprintf("%s (changed?) row %d col %d",
		em.Editor.GetFilename(), first.Line+1, first.Pos+1)
}
