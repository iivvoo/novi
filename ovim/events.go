package ovim

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

type Event interface{}

type KeyEvent struct {
	Modifier KeyModifier
	Key      KeyType
	Rune     rune
}

type CharacterEvent struct {
	Rune rune
}
