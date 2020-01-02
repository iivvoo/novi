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
 *
 * Could/should we support multiple cursors for vi emulation?
 * vim itself provides ctrl-v which is a bit like a multi-cursor, but not all command work on it
 *  (e.g. o or O have no effect. 'i' does have effecti, 'a' doesn't. Perhaps vim limitation?)
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
		Dispatch{Mode: ModeCommand, Events: []ovim.Event{
			ovim.CharacterEvent{Rune: 'i'},
			ovim.CharacterEvent{Rune: 'I'},
			ovim.CharacterEvent{Rune: 'o'},
			ovim.CharacterEvent{Rune: 'O'},
			ovim.CharacterEvent{Rune: 'a'},
			ovim.CharacterEvent{Rune: 'A'},
		}, Handler: em.HandleToModeEdit},

		Dispatch{Mode: ModeAny, Events: []ovim.Event{
			ovim.KeyEvent{Key: ovim.KeyLeft},
			ovim.KeyEvent{Key: ovim.KeyRight},
			ovim.KeyEvent{Key: ovim.KeyUp},
			ovim.KeyEvent{Key: ovim.KeyDown}}, Handler: em.HandleMoveCursors},
		Dispatch{Mode: ModeCommand, Events: []ovim.Event{
			ovim.CharacterEvent{Rune: 'h'},
			ovim.CharacterEvent{Rune: 'j'},
			ovim.CharacterEvent{Rune: 'k'},
			ovim.CharacterEvent{Rune: 'l'},
		}, Handler: em.HandleMoveHJKLCursors},
		Dispatch{Mode: ModeEdit, Event: ovim.CharacterEvent{}, Handler: em.HandleAnyRune},
	}
	em.dispatch = dispatch
	return em
}

// HandleToModeEdit handles the different switches to insert mode
func (em *Vi) HandleToModeEdit(ev ovim.Event) {
	em.Mode = ModeEdit

	r := ev.(ovim.CharacterEvent).Rune
	first := em.Editor.Cursors[0]

	switch r {
	case 'i':
		break // just insert at current cursor position
	case 'I':
		// insert at beginning of line
		em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorBegin)
	case 'o':
		// add line below current line
		// XXX TODO preserve indent (depend on indent mode?)
		em.Editor.Buffer.InsertLine(first, "", false)
		em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorDown)
	case 'O':
		// add line above cursor
		// XXX TODO preserve indent (depend on indent mode?)
		em.Editor.Buffer.InsertLine(first, "", true)
		// The cursor will already be at the inserted line, but may need to move to the start
		em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorBegin)
	case 'a':
		// after cursor
		em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorRight)
	case 'A':
		// at end
		em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorEnd)
	}
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
		if d.Do(event, em.Mode) {
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
