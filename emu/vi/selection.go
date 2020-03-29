package viemu

import "github.com/iivvoo/novi/novi"

// StartSelection initializes start/end with the current cursor position
func (em *Vi) StartSelection() {
	em.SelectionStart = *em.Editor.Cursors[0]
	em.Mode = ModeSelect
	em.Editor.Selection.Enable()
	em.UpdateSelection()
}

// CancelSelection cancels the selection, returns to Command mode
func (em *Vi) CancelSelection() {
	em.Selection = SelectionNone
	em.Mode = ModeCommand
	em.Editor.Selection.Disable()
}

// HandleCancelSelect is invoked when Escape is hit during selection
func (em *Vi) HandleCancelSelect(novi.Event) bool {
	em.CancelSelection()
	return true
}

// HandleSelectRemove handles selection removal keys, xdD
func (em *Vi) HandleSelectRemove(novi.Event) bool {
	// Block works differently of course
	s, e := em.GetEmuSelection()
	em.Editor.Buffer.RemoveBetweenCursors(&s, &e)
	em.CancelSelection()
	// move cursor to s
	c := em.Editor.Cursors[0]
	c.Line, c.Pos = s.Line, s.Pos
	return true
}

// HandleSelectChange handles selection change keys, cC
func (em *Vi) HandleSelectChange(novi.Event) bool {
	s, e := em.GetEmuSelection()
	em.Editor.Buffer.RemoveBetweenCursors(&s, &e)
	em.CancelSelection()
	em.Mode = ModeEdit
	c := em.Editor.Cursors[0]
	c.Line, c.Pos = s.Line, s.Pos
	return true
}

// GetEmuSelection translates the actual selection to how the emulation interprets them
func (em *Vi) GetEmuSelection() (novi.Cursor, novi.Cursor) {
	s, e := em.SelectionStart, em.SelectionEnd
	if e.Line < s.Line || (e.Line == s.Line && e.Pos < s.Pos) {
		// swap start, end
		s, e = e, s
	}

	switch em.Selection {
	case SelectionBlock:
		if e.Pos < s.Pos {
			e.Pos, s.Pos = s.Pos, e.Pos
		}
	case SelectionLines:
		s.Pos = 0
		e.Pos = em.Editor.Buffer.GetLine(e.Line).Len() - 1
	}
	return s, e
}

// UpdateSelection updates the end of the selection
func (em *Vi) UpdateSelection() {
	if em.Selection != SelectionNone {
		// There's a difference between the emulation selection and the UI selection.
		// In order to properly switch between the different types of selection, we need
		// to properly preserve the actual start/end, and only set the desired selection
		// (line,block,fluid) on em.Selection XXX
		// we have start/emuStart on the selection for that.

		em.SelectionEnd = *em.Editor.Cursors[0]

		em.Editor.Selection.SetBlock(em.Selection == SelectionBlock)
		s, e := em.GetEmuSelection()

		em.Editor.Selection.SetStart(s)
		em.Editor.Selection.SetEnd(e)
		log.Printf("Selection %s", em.Editor.Selection.ToString())
	}
}

// HandleSelectionBlock handles the block select key
func (em *Vi) HandleSelectionBlock(ev novi.Event) bool {
	if em.Selection == SelectionNone {
		em.Selection = SelectionBlock
		em.StartSelection()
	} else {
		em.Selection = SelectionBlock
		em.UpdateSelection()
	}
	return true
}

// HandleSelectionFluid handles the fluid select key
func (em *Vi) HandleSelectionFluid(ev novi.Event) bool {
	if em.Selection == SelectionNone {
		em.Selection = SelectionFluid
		em.StartSelection()
	} else {
		em.Selection = SelectionFluid
		em.UpdateSelection()
	}
	return true
}

// HandleSelectionLines handles the block select key
func (em *Vi) HandleSelectionLines(ev novi.Event) bool {
	if em.Selection == SelectionNone {
		em.Selection = SelectionLines
		em.StartSelection()
	} else {
		em.Selection = SelectionLines
		em.UpdateSelection()
	}
	return true
}
