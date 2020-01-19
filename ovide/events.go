package ovide

type Event interface{}

type QuitEvent struct{}

type OpenFileEvent struct {
	Filename string
	FullPath string
}
