package main

import (
	"fmt"
	"os"

	basicemu "gitlab.com/iivvoo/ovim/emu/basic"
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

	emu := basicemu.NewBasic(editor)
	core := ovim.NewCore(editor, ui, emu)
	core.Loop()
	ui.Finish()
}

func main() {
	start()
}
