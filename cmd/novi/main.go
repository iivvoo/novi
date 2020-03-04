package main

import (
	"flag"
	"strconv"
	"strings"
	"time"

	basicemu "github.com/iivvoo/novi/emu/basic"
	viemu "github.com/iivvoo/novi/emu/vi"
	"github.com/iivvoo/novi/logger"
	"github.com/iivvoo/novi/novi"
	termui "github.com/iivvoo/novi/ui/term"
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
	logger.OpenLog("novi.log")
	log.Printf("Starting at %s\n", time.Now())
	defer logger.CloseLog()

	editor := novi.NewEditor()
	ui := termui.NewTermUI(editor)
	ui.SetSize(w, h)
	defer novi.RecoverFromPanic(func() {
		ui.Finish()
	})

	fileName := ""

	if len(flag.Args()) > 0 {
		fileName = flag.Args()[0]
	}

	editor.LoadFile(fileName)
	editor.SetCursor(8, 0)

	var emu novi.Emulation

	if emuFlag == "vi" {
		emu = viemu.NewVi(editor)
	} else {
		emu = basicemu.NewBasic(editor)
	}
	core := novi.NewCore(editor, ui, emu)
	core.Loop()
	ui.Finish()
}

func main() {
	start()
}
