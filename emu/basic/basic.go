package basicemu

import (
	"fmt"

	"github.com/iivvoo/ovim/logger"
	"github.com/iivvoo/ovim/ovim"
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
	Editor *ovim.Editor
}

func NewBasic(e *ovim.Editor) *Basic {
	return &Basic{Editor: e}
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
			Move(c, ovim.CursorLeft)
		} else if c.Line > 0 {
			// first move the cursor so we can use CursorEnd to move to the desired position
			l := c.Line
			Move(c, ovim.CursorUp)
			Move(c, ovim.CursorEnd)
			em.Editor.Buffer.JoinLineWithPrevious(l)

			// adjust all other cursors that are on/after l
			// XXX Untested
			// XXX also wrong, also changes cursors *before* line.
			for _, cc := range em.Editor.Cursors.After(c) {
				Move(cc, ovim.CursorUp)
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
func (em *Basic) HandleEvent(event ovim.Event) bool {
	switch ev := event.(type) {
	case *ovim.KeyEvent:
		// control keys, purely control
		if ev.Modifier == ovim.ModCtrl {
			switch ev.Rune {
			case 'h':
				em.Backspace()
			case 'q':
				return false
			case 's':
				em.Editor.SaveFile()
				log.Println("File saved")
			default:
				log.Printf("Don't know what to do with control key %+v %c", ev, ev.Rune)
			}
			// no modifier at all
		} else if ev.Modifier == 0 {
			switch ev.Key {
			case ovim.KeyBackspace, ovim.KeyDelete:
				// for cursors on pos 0, join with prev (if any)
				em.Backspace()
			case ovim.KeyEnter:
				for _, c := range em.Editor.Cursors {
					em.Editor.Buffer.SplitLine(c)
					Move(c, ovim.CursorDown)
					Move(c, ovim.CursorBegin)
					// update all cursors after
					for _, ca := range em.Editor.Cursors.After(c) {
						ca.Line++
					}
				}
			case ovim.KeyLeft, ovim.KeyRight, ovim.KeyUp, ovim.KeyDown, ovim.KeyHome, ovim.KeyEnd:
				for _, c := range em.Editor.Cursors {
					Move(c, ovim.CursorMap[ev.Key])
				}
			default:
				log.Printf("Don't know what to do with key event %+v", ev)
			}
		}
	case *ovim.CharacterEvent:
		em.Editor.Buffer.PutRuneAtCursors(em.Editor.Cursors, ev.Rune)
		for _, c := range em.Editor.Cursors {
			Move(c, ovim.CursorRight)
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
