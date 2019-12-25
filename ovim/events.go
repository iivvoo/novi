package ovim

type KeyType uint16

const (
	KeyEscape KeyType = iota
	KeyEnter
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
)

type Event interface{}

type KeyEvent struct {
	Key KeyType
}

type CharacterEvent struct {
	Rune rune
}
