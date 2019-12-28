package main

import (
	"fmt"
	"os"
	"time"

	basicemu "gitlab.com/iivvoo/ovim/emu/basic"
	"gitlab.com/iivvoo/ovim/logger"
	"gitlab.com/iivvoo/ovim/ovim"
	termui "gitlab.com/iivvoo/ovim/ui/term"
)

var log = logger.GetLogger("main")

func start() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ovim filename")
		os.Exit(1)
	}
	// initialize logger
	logger.OpenLog("ovim.log")
	log.Printf("Starting at %s\n", time.Now())
	defer logger.CloseLog()

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
