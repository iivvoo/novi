package viemu

import (
	"fmt"
	"strings"

	"github.com/iivvoo/ovim/logger"
	"github.com/iivvoo/ovim/ovim"
)

/*
 *
 * Lots of stuff to do. Start with basic non-ex (?) commands, controls:
 * insert: iIoOaA OK (single cursor)
 * D - delete (+ insert) to EOL
 * C - change to EOL
 * regular <enter> in insert mode OK
 * HML - high / middle / low (screen, niet file)
 * vi scrolls few lines before top/bottom, not at
 * <num?>gg (top) G (end) of file OK
 * w (jump word), with counter. Keep support for "c<n>w" in mind!
 *  w = non-space? W=true word?
 *  c<N>w and d<N>w are simiar (c is d with insert)
 * ZQ to exit without save OK
 * copy/paste (non/term/mouse: y, p etc)
 * commands: d10d, c5w, 10x, etc.
 * proper tab support
 * ^R" (paste buffer magic) in insert mode
 *
 * Could/should we support multiple cursors for vi emulation?
 * vim itself provides ctrl-v which is a bit like a multi-cursor, but not all command work on it
 *  (e.g. o or O have no effect. 'i' does have effecti, 'a' doesn't. Perhaps vim limitation?)
 *
 * '.' replays last command - we need a way to "store" this (is storing the keypresses sufficient?)
 *
 */

var log = logger.GetLogger("viemu")

// ViMode is the current mode of operation
type ViMode int

// It currently has these modes
const (
	ModeAny ViMode = iota
	ModeEdit
	ModeCommand
)

// DispatchHandler is the signature for a handler in the dispatch table
type DispatchHandler func(ovim.Event) bool

// Dispatch maps a Key/CharacterEvent to a handler
type Dispatch struct {
	Mode    ViMode
	Event   ovim.Event
	Events  []ovim.Event
	Handler DispatchHandler
}

// Do calls the handler if the event matches
func (d Dispatch) Do(event ovim.Event, mode ViMode) bool {
	if event.Equals(d.Event) && (d.Mode == ModeAny || d.Mode == mode) {
		return d.Handler(event)
	}
	for _, e := range d.Events {
		if event.Equals(e) && (d.Mode == ModeAny || d.Mode == mode) {
			return d.Handler(e)
		}
	}
	return false
}

// Vi encapsulate all the Vi emulation state
type Vi struct {
	Editor        *ovim.Editor
	Mode          ViMode
	CommandBuffer string
	Counter       int

	dispatch []Dispatch
}

/*
 * in command mode, you don't simply press a key. It can be prefixed with a count
 * l = right
 * 10l = 10 right (as far as possible)
 * x = remove char
 * 10x = remove 10 chars
 *
 * 'd' by itself is nothing, it's always 'dd' which can be 10dd or d10d. 2d2d also works -> 4dd
 * actually 2d2d is 2*(d2d), so 3d2d will delete 6, not 5
 *
 * Certain commands will clear any counter and just work, e.g. the insertion keys
 * Escape in command mode clears command buffer
 *
 * Approach: put everything in a buffer. After each key, check of buffer is a complete command
 *
 * Vim extra's:
 * insert keys can be commands/repeated: 3iYes 3oHello
 * e/E go to end of word, similar to w/W
 * 'r' replace single character under cursor, can also be used with count 5ra -> replace with aaaaa
 * c|d<n>w werkt net iets anders dan regulier w - verandert tot voor matchend woord, whitespace in tact
 */

// NewVi creates/setups up a new Vi emulation instance
func NewVi(e *ovim.Editor) *Vi {
	em := &Vi{Editor: e, Mode: ModeCommand}
	dispatch := []Dispatch{
		Dispatch{Mode: ModeEdit, Event: ovim.KeyEvent{Key: ovim.KeyEscape}, Handler: em.HandleToModeCommand},
		Dispatch{Mode: ModeCommand, Event: ovim.KeyEvent{Key: ovim.KeyEscape}, Handler: em.HandleCommandClear},
		Dispatch{Mode: ModeCommand, Event: ovim.KeyEvent{Key: ovim.KeyEnter}, Handler: em.HandleCommandEnter},
		Dispatch{Mode: ModeEdit, Event: ovim.KeyEvent{Key: ovim.KeyEnter}, Handler: em.HandleEditEnter},
		Dispatch{Mode: ModeCommand, Events: []ovim.Event{
			ovim.CharacterEvent{Rune: 'i'},
			ovim.CharacterEvent{Rune: 'I'},
			ovim.CharacterEvent{Rune: 'o'},
			ovim.CharacterEvent{Rune: 'O'},
			ovim.CharacterEvent{Rune: 'a'},
			ovim.CharacterEvent{Rune: 'A'},
		}, Handler: em.HandleInsertionKeys},

		Dispatch{Mode: ModeAny, Events: []ovim.Event{
			ovim.KeyEvent{Key: ovim.KeyLeft},
			ovim.KeyEvent{Key: ovim.KeyRight},
			ovim.KeyEvent{Key: ovim.KeyUp},
			ovim.KeyEvent{Key: ovim.KeyDown},
			ovim.KeyEvent{Key: ovim.KeyEnd},
			ovim.KeyEvent{Key: ovim.KeyHome},
		}, Handler: em.HandleMoveCursors},
		Dispatch{Mode: ModeAny, Events: []ovim.Event{
			ovim.KeyEvent{Key: ovim.KeyBackspace},
			ovim.KeyEvent{Key: ovim.KeyDelete},
		}, Handler: em.HandleBackspace},
		// Sort of a generic fallthrough handler - handles commands in command mode
		Dispatch{Mode: ModeCommand, Event: ovim.CharacterEvent{}, Handler: em.HandleCommandBuffer},
		Dispatch{Mode: ModeEdit, Event: ovim.CharacterEvent{}, Handler: em.HandleAnyRune},
	}
	em.dispatch = dispatch
	return em
}

// HandleCommandEnter handles enter in command mode
func (em *Vi) HandleCommandEnter(ev ovim.Event) bool {
	for _, c := range em.Editor.Cursors {
		Move(c, ovim.CursorDown)
		Move(c, ovim.CursorBegin)
	}
	return true
}

// HandleEditEnter handles enter in insert mode
func (em *Vi) HandleEditEnter(ev ovim.Event) bool {
	// XXX identical to "basic" emulation
	for _, c := range em.Editor.Cursors {
		em.Editor.Buffer.SplitLine(c)
		Move(c, ovim.CursorDown)
		Move(c, ovim.CursorBegin)
		// update all cursors after
		for _, ca := range em.Editor.Cursors.After(c) {
			ca.Line++
		}
	}
	return true
}

// HandleBackspace handles backspace behaviour in both edit and command mode
func (em *Vi) HandleBackspace(ev ovim.Event) bool {
	// BUG: vim seems to allow counts on backspace in commandmode,
	// e.g. 10<backspace>, so it should be handled there
	if em.Mode == ModeCommand {
		for _, c := range em.Editor.Cursors {
			if c.Pos == 0 && c.Line != 0 {
				Move(c, ovim.CursorUp)
				Move(c, ovim.CursorEnd)
			} else {
				Move(c, ovim.CursorLeft)
			}
		}
	} else {
		for _, c := range em.Editor.Cursors {
			if c.Pos > 0 {
				em.Editor.Buffer.RemoveRuneBeforeCursor(c)
				Move(c, ovim.CursorLeft)
			} else {
				// identical to basic emulation XXX
				l := c.Line
				Move(c, ovim.CursorUp)
				Move(c, ovim.CursorEnd)
				em.Editor.Buffer.JoinLineWithPrevious(l)
				// except here, since "End" in vi moves to the last character, not past it, for which we need to compensate
				Move(c, ovim.CursorRight)

				for _, cc := range em.Editor.Cursors.After(c) {
					Move(cc, ovim.CursorUp)
				}
			}
		}
	}
	return true
}

// HandleCommandClear clears the current command state (if any)
func (em *Vi) HandleCommandClear(ev ovim.Event) bool {
	em.CommandBuffer = ""
	return true
}

// RemoveCharacters removes a number of characters before or after the cursors
func (em *Vi) RemoveCharacters(howmany int, before bool) {
	for _, c := range em.Editor.Cursors {
		em.Editor.Buffer.RemoveCharacters(c, before, howmany)
		if before {
			MoveMany(c, ovim.CursorLeft, howmany)
		}
	}
}

// RemoveLines removes full lines
func (em *Vi) RemoveLines(howmany int) {
	// How would this behave on multiple cursors?
	first := em.Editor.Cursors[0]
	for i := 0; i < howmany; i++ {
		if !em.Editor.Buffer.RemoveLine(first.Line) {
			// We ran out of lines, no need to continue, but do move up
			Move(first, ovim.CursorUp)
			break
		}
	}
	first.Validate()
}

// JumpStartEndLine handles jumping to the start/end of line
func (em *Vi) JumpStartEndLine(howmany int, jumpstart bool) {
	for _, c := range em.Editor.Cursors {
		if jumpstart {
			// howmany has no meaning
			c.Pos = 0
		} else {
			MoveMany(c, ovim.CursorDown, howmany-1)
			Move(c, ovim.CursorEnd)
		}
	}
}

// HandleEvent is the main entry point
func (em *Vi) HandleEvent(event ovim.Event) bool {
	for _, d := range em.dispatch {
		if d.Do(event, em.Mode) {
			// returns false if we need to exit
			return em.CheckExecuteCommandBuffer()
		}
	}
	return false
}

// HandleInsertionKeys handles the different switches to insert mode
func (em *Vi) HandleInsertionKeys(ev ovim.Event) bool {
	em.Mode = ModeEdit

	r := ev.(ovim.CharacterEvent).Rune
	first := em.Editor.Cursors[0]

	switch r {
	case 'i': // just insert at current cursor position
		break
	case 'I': // insert at beginning of line
		Move(first, ovim.CursorBegin)
	case 'o': // add line below current line
		// XXX TODO preserve indent (depend on indent mode?)
		em.Editor.Buffer.InsertLine(first, "", false)
		Move(first, ovim.CursorDown)
	case 'O': // add line above cursor
		// XXX TODO preserve indent (depend on indent mode?)
		em.Editor.Buffer.InsertLine(first, "", true)
		// The cursor will already be at the inserted line, but may need to move to the start
		Move(first, ovim.CursorBegin)
	case 'a': // after cursor
		Move(first, ovim.CursorRight)
	case 'A': // at end
		Move(first, ovim.CursorEnd)
		first.Pos++
	}
	return true
}

// HandleToModeCommand simply switches (back) to command mode
func (em *Vi) HandleToModeCommand(ovim.Event) bool {
	em.Mode = ModeCommand
	return true
}

// HandleAnyRune simply inserts the character in edit mode
func (em *Vi) HandleAnyRune(ev ovim.Event) bool {
	r := ev.(*ovim.CharacterEvent).Rune
	em.Editor.Buffer.PutRuneAtCursors(em.Editor.Cursors, r)
	for _, c := range em.Editor.Cursors {
		// Move(CursorRight) won't do since it will restrict to the last character
		c.Pos++
	}
	return true
}

// HandleCommandBuffer handles all keys that affect the command buffer
func (em *Vi) HandleCommandBuffer(ev ovim.Event) bool {
	commands := "BbdgGhjklxXdwWcZQ0123456789$^"
	r := ev.(*ovim.CharacterEvent).Rune

	if strings.IndexRune(commands, r) != -1 {
		em.CommandBuffer += string(r)
		return true
	}
	return false
}

// SaveFile tells the editor to save the file
func (em *Vi) SaveFile() {
	em.Editor.SaveFile()
}

// CheckExecuteCommandBuffer checks if there's a full, complete command and, if so, executes it
func (em *Vi) CheckExecuteCommandBuffer() bool {
	/*
	 * a vi(m?) command has the structure
	 * <number?>character
	 * <number?>character(<number?>character)? e.g. 2d3d -> 6dd, or d10d -> 10dd
	 *
	 * (vim actually understands <num><keyup>!, same for backspace)
	 *
	 * "just" 0 = Begin of line
	 * odd case, 2d0 deletes current line to beginning
	 *
	 * There are also combinations, e.g c3w -> what about 2c3w?
	 */

	count, command := ParseCommand(em.CommandBuffer)
	switch command {
	case "h", "j", "k", "l":
		em.MoveCursorRune(rune(command[0]), count)
		em.CommandBuffer = ""
	case "w", "W", "b", "B":
		em.JumpWord(rune(command[0]), count)
		em.CommandBuffer = ""
	case "x", "X":
		em.RemoveCharacters(count, command == "X")
		em.CommandBuffer = ""
	case "gg", "G":
		em.JumpTopBottom(count, command == "gg")
		em.CommandBuffer = ""
	case "^", "$":
		em.JumpStartEndLine(count, command == "^")
		em.CommandBuffer = ""
	case "ZZ":
		em.SaveFile()
		em.CommandBuffer = ""
	case "ZQ":
		em.CommandBuffer = ""
		return false // signals exit
	case "dd":
		em.RemoveLines(count)
		em.CommandBuffer = ""
	}
	return true
}

// JumpWord jumps to the next word / sequence of separators
func (em *Vi) JumpWord(r rune, howmany int) {
	for i := 0; i < howmany; i++ {
		for _, c := range em.Editor.Cursors {
			switch r {
			case 'W':
				l, p := JumpWordForward(em.Editor.Buffer, c)
				c.Line, c.Pos = l, p
			case 'w':
				l, p := JumpForward(em.Editor.Buffer, c)
				c.Line, c.Pos = l, p
			case 'B':
				l, p := JumpWordBackward(em.Editor.Buffer, c)
				c.Line, c.Pos = l, p
			case 'b':
				l, p := JumpBackward(em.Editor.Buffer, c)
				c.Line, c.Pos = l, p
			}
		}
	}

}

// JumpTopBottom handles jumping using the gg / G command
func (em *Vi) JumpTopBottom(howmany int, jumptop bool) {
	// if howany is > 1, it's always a jump from the top
	for _, c := range em.Editor.Cursors {
		if howmany > 1 {
			c.Line = 0
			c.Pos = 0
			MoveMany(c, ovim.CursorDown, howmany-1)
		} else if jumptop {
			// this will move all cursors to (0,0) -- remove them?
			c.Line = 0
			c.Pos = 0
		} else {
			c.Line = em.Editor.Buffer.Length() - 1
			c.Pos = 0
		}
	}
}

// GetStatus provides a way for the Editor to get the emulation's status
func (em *Vi) GetStatus(width int) string {
	mode := ""
	modified := ""
	first := em.Editor.Cursors[0]
	if em.Mode == ModeEdit {
		mode = "--INSERT-- "
	}
	if em.Editor.Buffer.Modified {
		modified = "(modified) "
	}
	return mode + fmt.Sprintf("%s %s   %s  row %d col %d",
		em.Editor.GetFilename(), modified, em.CommandBuffer, first.Line+1, first.Pos+1)
}
