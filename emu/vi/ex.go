package viemu

import "github.com/iivvoo/ovim/ovim"

type Input struct {
	Buffer *Line
	Pos    int
}

func NewInput() *Input {
	return &Input{NewLine(), 0}
}

func (i *Input) Clear() {
	i.Pos = 0
	i.Buffer = NewLine()
}

func (i *Input) Len() int {
	return i.Buffer.Len()
}

func (i *Input) Backspace() {
	if i.Pos > 0 {
		i.Buffer.RemoveRune(i.Pos - 1)
		i.Pos--
	}
}

func (i *Input) CursorLeft() {
	if i.Pos > 0 {
		i.Pos--
	}
}

func (i *Input) CursorRight() {
	if i.Pos < i.Buffer.Len() {
		i.Pos++
	}
}

func (i *Input) Insert(r rune) {
	i.Buffer.InsertRune(r, i.Pos)
	i.Pos++
}

// Ex encapsulates the state of the ex buffer / mode
type Ex struct {
	input *Input
}

func NewEx() *Ex {
	return &Ex{input: NewInput()}
}

func (ex *Ex) Clear() {
	ex.input.Clear()
}

/*
 * TODO:
 * minimal buffer editing (cursor, backspace)
 * backspace on empty = cancel
 *
 * Eventually we'll need similar input for /search, so reuse is desirable
 */

// HandleExCommand handles the ':' ex commands
func (em *Vi) HandleExCommand() {
	/*
	 * could scan the exBuffer continuously and adjust the buffer, e.g. highlight matches. But for now,
	 * just handle commands such as:
	 *
	 * :q :q!
	 * :w :w!
	 * :<linenumber>
	 */
}

// HandleExKey handles non-character "special" keys such as cursor keys, escape, backspace
func (em *Vi) HandleExKey(e *ovim.KeyEvent) {
	/*
	   Left/right: move cursor
	   backspace: remove characters, escape if empty

	   nice to have:
	   up/down> history (if empty)
	*/
	switch e.Key {
	case ovim.KeyBackspace:
		if em.ex.input.Len() == 0 {
			em.c <- &ovim.CloseInputEvent{ID: 1}
		}
		em.ex.input.Backspace()
	case ovim.KeyLeft:
		em.ex.input.CursorLeft()
	case ovim.KeyRight:
		em.ex.input.CursorRight()
	case ovim.KeyEscape:
		em.c <- &ovim.CloseInputEvent{ID: 1}
	case ovim.KeyEnter:
		// handle ex command
		em.c <- &ovim.CloseInputEvent{ID: 1}
	}
	em.c <- &ovim.UpdateInputEvent{ID: 1, Text: em.ex.input.Buffer.ToString(), Pos: em.ex.input.Pos}
}

// HandleExInput handles the Ex input events
func (em *Vi) HandleExInput(event ovim.Event) bool {
	if char, ok := event.(*ovim.CharacterEvent); ok {
		em.ex.input.Insert(char.Rune)
		em.c <- &ovim.UpdateInputEvent{ID: 1, Text: em.ex.input.Buffer.ToString(), Pos: em.ex.input.Pos}
	} else if key, ok := event.(*ovim.KeyEvent); ok {
		em.HandleExKey(key)
	}
	return true
}
