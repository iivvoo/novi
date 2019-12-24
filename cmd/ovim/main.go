package main

import (
	"fmt"
	"os"

	"gitlab.com/iivvoo/ovim/ovim"
)

func start() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ovim filename")
		os.Exit(1)
	}

	editor := ovim.NewEditor()
	ui := ovim.NewTermUI(editor)
	defer ovim.RecoverFromPanic(func() {
		ui.Finish()
	})

	fileName := os.Args[1]

	editor.LoadFile(fileName)
	editor.SetCursor(8, 0)

	ui.Loop()
}

func main() {
	start()
}
