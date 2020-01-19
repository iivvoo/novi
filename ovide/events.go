package ovide

type Event interface{}

type DebugEvent struct {
	Msg string
}

type QuitEvent struct{}

type OpenFileEvent struct {
	Filename string
	FullPath string
}

type NewFileEvent struct {
	ParentFolder string
}

type NewFolderEvent struct {
	ParentFolder string
}
