package main

import (
	"flag"
	"strconv"
	"strings"
	"time"

	basicemu "github.com/iivvoo/ovim/emu/basic"
	viemu "github.com/iivvoo/ovim/emu/vi"
	"github.com/iivvoo/ovim/logger"
	"github.com/iivvoo/ovim/ovim"
	termui "github.com/iivvoo/ovim/ui/term"
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

	fileName := ""

	if len(flag.Args()) > 0 {
		fileName = flag.Args()[0]
	}

	editor.LoadFile(fileName)
	editor.SetCursor(8, 0)

	var emu ovim.Emulation

	if emuFlag == "vi" {
		emu = viemu.NewVi(editor)
	} else {
		emu = basicemu.NewBasic(editor)
	}
	core := ovim.NewCore(editor, ui, ui, ui, emu)
	core.Loop()
	ui.Finish()
}

func main() {
	start()
}
