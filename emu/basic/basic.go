package basicemu

import (
	"fmt"

	"github.com/iivvoo/novi/logger"
	"github.com/iivvoo/novi/novi"
)

var log = logger.GetLogger("basicemu")

/*
 * Emulate a regular, basic editor. Standard keybinding/behaviour
 */

/*
 * In general, an editor emulator is responsible for handling all/most keys.
 * It should probbaly also expose which keys it wants to handle, so the overall
 * framework can decide what keys are available for additional bindings
 *
 * Ctrl-s - save
 * Ctrl-q - quite
 * Home/End - start/end line
 * pgup, pgdn
 * insert: toggle overwrite/insert-move (include change cursor)
 */

// The Basic struct encapsulates all state for the Basic Editing emulation
type Basic struct {
	Editor *novi.Editor

	c chan novi.EmuEvent
}

func NewBasic(e *novi.Editor) *Basic {
	return &Basic{Editor: e}
}

// SetChan passes us a channel to communicate with the core.
func (em *Basic) SetChan(c chan novi.EmuEvent) {
	em.c = c
}

/*
 * The emulation need to interact directly with the editor (and possibly UI, Core)
 * so no loop/channel magic.
 *
 * We do need to signal quit to the editor.
 * Furthermore,
 * - emulation may do complex tasks such as saving, asking filename
 * - emulation may need to know UI specifics, e.g. screenheight for pgup/pgdn
 */

func (em *Basic) Backspace() {
	for _, c := range em.Editor.Cursors {
		if c.Pos > 0 {
			em.Editor.Buffer.RemoveRuneBeforeCursor(c)
			Move(c, novi.CursorLeft)
		} else if c.Line > 0 {
			// first move the cursor so we can use CursorEnd to move to the desired position
			l := c.Line
			Move(c, novi.CursorUp)
			Move(c, novi.CursorEnd)
			em.Editor.Buffer.JoinLineWithPrevious(l)

			// adjust all other cursors that are on/after l
			// XXX Untested
			// XXX also wrong, also changes cursors *before* line.
			for _, cc := range em.Editor.Cursors.After(c) {
				Move(cc, novi.CursorUp)
			}
		}
	}
}

/*
 * How should we handle/manipulate lines, buffers?
 * - A bunch of methods on Editor, hoping they will satisfy al needs,
 *   with correct multi-cursor behaviours
 * - Direct manipulation of text buffer
 *
 * Also, who is in charge of updating the cursor(s)?
 */
func (em *Basic) HandleEvent(_ novi.InputID, event novi.Event) bool {
	switch ev := event.(type) {
	case *novi.KeyEvent:
		// control keys, purely control
		if ev.Modifier == novi.ModCtrl {
			switch ev.Rune {
			case 'h':
				em.Backspace()
			case 'q':
				return false
			case 's':
				em.c <- &novi.SaveEvent{}
				log.Println("File saved")
			default:
				log.Printf("Don't know what to do with control key %+v %c", ev, ev.Rune)
			}
			// no modifier at all
		} else if ev.Modifier == 0 {
			switch ev.Key {
			case novi.KeyBackspace, novi.KeyDelete:
				// for cursors on pos 0, join with prev (if any)
				em.Backspace()
			case novi.KeyEnter:
				for _, c := range em.Editor.Cursors {
					em.Editor.Buffer.SplitLine(c)
					Move(c, novi.CursorDown)
					Move(c, novi.CursorBegin)
					// update all cursors after
					for _, ca := range em.Editor.Cursors.After(c) {
						ca.Line++
					}
				}
			case novi.KeyLeft, novi.KeyRight, novi.KeyUp, novi.KeyDown, novi.KeyHome, novi.KeyEnd:
				for _, c := range em.Editor.Cursors {
					Move(c, novi.CursorMap[ev.Key])
				}
			default:
				log.Printf("Don't know what to do with key event %+v", ev)
			}
		}
	case *novi.CharacterEvent:
		em.Editor.Buffer.PutRuneAtCursors(em.Editor.Cursors, ev.Rune)
		for _, c := range em.Editor.Cursors {
			Move(c, novi.CursorRight)
		}
	default:
		log.Printf("Don't know what to do with event %+v", ev)
	}
	return true
}

// GetStatus returns the status for the emulation, to be printed in a status bar
func (em *Basic) GetStatus(width int) string {
	first := em.Editor.Cursors[0]
	changed := ""

	if em.Editor.Buffer.Modified {
		changed = "(changed) "
	}
	// Make use of width to align cursor row/col right. Truncate if necessary
	return fmt.Sprintf("%s %srow %d col %d (INS)",
		em.Editor.GetFilename(), changed, first.Line+1, first.Pos+1)
}
