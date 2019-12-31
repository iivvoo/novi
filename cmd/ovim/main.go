package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	basicemu "gitlab.com/iivvoo/ovim/emu/basic"
	viemu "gitlab.com/iivvoo/ovim/emu/vi"
	"gitlab.com/iivvoo/ovim/logger"
	"gitlab.com/iivvoo/ovim/ovim"
	termui "gitlab.com/iivvoo/ovim/ui/term"
)

var log = logger.GetLogger("main")

func start() {

	var sizeFlag string
	var emuFlag string
	var w, h int
	var err error

	flag.StringVar(&sizeFlag, "area", "", "Edit area size")
	flag.StringVar(&emuFlag, "emu", "basic", "Emulation to use")

	flag.Parse()
	if sizeFlag != "" {
		p := strings.SplitN(sizeFlag, "x", 2)
		if len(p) != 2 {
			panic("area")
		}
		w, err = strconv.Atoi(p[0])
		if err != nil {
			panic("area parse w")
		}
		h, err = strconv.Atoi(p[1])
		if err != nil {
			panic("area parse h")
		}
	}

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: ovim filename")
		os.Exit(1)
	}
	// initialize logger
	logger.OpenLog("ovim.log")
	log.Printf("Starting at %s\n", time.Now())
	defer logger.CloseLog()

	editor := ovim.NewEditor()
	ui := termui.NewTermUI(editor)
	ui.SetSize(w, h)
	defer ovim.RecoverFromPanic(func() {
		ui.Finish()
	})

	fileName := flag.Args()[0]

	editor.LoadFile(fileName)
	editor.SetCursor(8, 0)

	var emu ovim.Emulation

	if emuFlag == "vi" {
		emu = viemu.NewVi(editor)
	} else {
		emu = basicemu.NewBasic(editor)
	}
	core := ovim.NewCore(editor, ui, emu)
	core.Loop()
	ui.Finish()
}

func main() {
	start()
}
