package viemu

import (
	"strconv"
	"strings"

	"github.com/iivvoo/ovim/ovim"
)

// Input holds the state of a simple input buffer
type Input struct {
	Buffer *Line
	Pos    int
}

// NewInput creates a new input
func NewInput() *Input {
	return &Input{NewLine(), 0}
}

// Clear clears the input, resetting its size/contents to empty
func (i *Input) Clear() {
	i.Pos = 0
	i.Buffer = NewLine()
}

// ToString returns the string representation of the contents of the input
func (i *Input) ToString() string {
	return i.Buffer.ToString()
}

// Len returns the length of the input
func (i *Input) Len() int {
	return i.Buffer.Len()
}

// Backspace performs the backspace operation on the input on the current cursor position
func (i *Input) Backspace() {
	if i.Pos > 0 {
		i.Buffer.RemoveRune(i.Pos - 1)
		i.Pos--
	}
}

// CursorLeft moves the cursor left 1 position, if possible
func (i *Input) CursorLeft() {
	if i.Pos > 0 {
		i.Pos--
	}
}

// CursorRight moves the cursor right 1 position, if possible
func (i *Input) CursorRight() {
	if i.Pos < i.Buffer.Len() {
		i.Pos++
	}
}

// Insert inserts a rune at the current cursor position and advances the cursor
func (i *Input) Insert(r rune) {
	i.Buffer.InsertRune(r, i.Pos)
	i.Pos++
}

// Ex encapsulates the state of the ex buffer / mode
type Ex struct {
	input *Input
}

// NewEx creates a new Ex instance
func NewEx() *Ex {
	return &Ex{input: NewInput()}
}

// Clear clears the ex instance (input)
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

	cmd := strings.TrimSpace(em.ex.input.ToString())

	if cmd == "" {
		return
	}

	if number, err := strconv.Atoi(cmd); err == nil {
		if number <= 0 {
			number = 1
		}
		if number > em.Editor.Buffer.Length() {
			number = em.Editor.Buffer.Length()
		}
		first := em.Editor.Cursors[0]
		first.Line = number - 1
		return
	}

	parts := strings.Split(cmd, " ")
	p := parts[0]
	l := len(parts)
	// may contain filename?
	switch p {
	case "$": // last line of file
		// Error if any additional arguments
		if l > 1 {
			em.c <- &ovim.ErrorEvent{Message: "Extra characters after command"}
			return
		}
		em.JumpTopBottom(0, false)
	case "w", "wq", "w!", "wq!":
		if l > 2 {
			em.c <- &ovim.ErrorEvent{Message: "Extra characters after command"}
			return
		}
		force := strings.ContainsRune(p, '!')
		quit := strings.ContainsRune(p, 'q')
		fname := ""
		if l > 1 {
			fname = parts[1]
		}
		em.c <- &ovim.SaveEvent{Name: fname, Force: force}
		if quit {
			em.c <- &ovim.QuitEvent{Force: force}
		}
	case "q", "q!":
		if l > 1 {
			em.c <- &ovim.ErrorEvent{Message: "Extra characters after command"}
			return
		}
		force := strings.ContainsRune(p, '!')
		em.c <- &ovim.QuitEvent{Force: force}
	}
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
		em.ex.Clear()
		em.c <- &ovim.CloseInputEvent{ID: 1}
		return
	case ovim.KeyEnter:
		em.HandleExCommand()
		em.ex.Clear()
		em.c <- &ovim.CloseInputEvent{ID: 1}
		return
	}
	em.c <- &ovim.UpdateInputEvent{ID: 1, Text: em.ex.input.ToString(), Pos: em.ex.input.Pos}
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
