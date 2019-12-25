package main

import (
	"fmt"
	"os"

	"gitlab.com/iivvoo/ovim/ovim"
	termui "gitlab.com/iivvoo/ovim/ui/term"
)

func start() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ovim filename")
		os.Exit(1)
	}

	editor := ovim.NewEditor()
	ui := termui.NewTermUI(editor)
	defer ovim.RecoverFromPanic(func() {
		ui.Finish()
	})

	fileName := os.Args[1]

	editor.LoadFile(fileName)
	editor.SetCursor(8, 0)

	core := ovim.NewCore(editor, ui, nil)
	core.Loop()
	ui.Finish()
}

func main() {
	start()
}
