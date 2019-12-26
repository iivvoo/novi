package basicemu

import (
	"fmt"

	"gitlab.com/iivvoo/ovim/ovim"
)

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
			case 'q':
				return false
			default:
				panic(fmt.Sprintf("Unknown ctrl key: %c", ev.Rune))
			}
			// no modifier at all
		} else if ev.Modifier == 0 {
			switch ev.Key {
			case ovim.KeyBackspace:

			case ovim.KeyEnter:
				em.Editor.Buffer.AddLine(ovim.Line(""))
			case ovim.KeyLeft, ovim.KeyRight, ovim.KeyUp, ovim.KeyDown:
				em.Editor.MoveCursor(ev.Key)
			default:
				// write unhandled key to log / statusbar
				panic(ev)
			}
		}
	case *ovim.CharacterEvent:
		em.Editor.Buffer.PutRuneAtCursors(em.Editor.Cursors, ev.Rune)
		// passing a key here is weird
		em.Editor.Cursors.Move(em.Editor.Buffer, ovim.CursorRight)

		// update cursors
	default:
		panic(ev)
	}
	return true
}
