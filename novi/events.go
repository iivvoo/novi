package novi

type KeyModifier uint8

const (
	ModShift KeyModifier = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
)

type KeyType uint16

const (
	KeyRune KeyType = iota
	KeyEscape
	KeyEnter
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDn
	KeyBackspace
	KeyTab
	KeyDelete
	KeyInsert
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	/*
		tcell defines these and many others, but limit to most common keys for now
			KeyClear:          "Clear",
			KeyExit:           "Exit",
			KeyCancel:         "Cancel",
			KeyPause:          "Pause",
			KeyPrint:          "Print",
	*/
)

type InputSource int
type Event interface { // UIEvent?
	Equals(Event) bool
	SetSource(InputSource)
	GetSource() InputSource
}

type KeyEvent struct {
	Modifier KeyModifier
	Key      KeyType
	Rune     rune
	Source   InputSource
}

func (e *KeyEvent) Equals(other Event) bool {
	if ke, ok := other.(*KeyEvent); ok {
		return ke.Source == e.Source && ke.Modifier == e.Modifier && ke.Key == e.Key && (ke.Rune == 0 || ke.Rune == e.Rune)
	}
	return false
}

func (e *KeyEvent) SetSource(s InputSource) {
	e.Source = s
}

func (e *KeyEvent) GetSource() InputSource {
	return e.Source
}

type CharacterEvent struct {
	Rune   rune
	Source InputSource
}

func (e *CharacterEvent) Equals(other Event) bool {
	if ce, ok := other.(*CharacterEvent); ok {
		return ce.Source == e.Source && (ce.Rune == 0 || ce.Rune == e.Rune)
	}
	return false
}

func (e *CharacterEvent) SetSource(s InputSource) {
	e.Source = s
}

func (e *CharacterEvent) GetSource() InputSource {
	return e.Source
}

type InputID int

// This is/should be a different kind of event
type EmuEvent interface{}

type AskInputEvent struct {
	ID     InputID
	Prompt string
}
type CloseInputEvent struct {
	ID InputID
}

type UpdateInputEvent struct {
	ID   InputID
	Text string
	Pos  int
}

type QuitEvent struct {
	Force bool
}

type SaveEvent struct {
	Name  string
	Force bool
}

type ErrorEvent struct {
	Message string
}
