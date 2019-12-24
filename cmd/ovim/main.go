package main

import (
	"gitlab.com/iivvoo/ovim/ovim"
)

func start() {

	editor := ovim.NewEditor()
	ui := ovim.NewTermUI(editor)
	defer ovim.RecoverFromPanic(func() {
		ui.Finish()
	})

	text := []string{
		"Hello world",
		"",
		"So nice you see you",
		"this is a reeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeally long line",
		"bye now",
		"",
		"",
		"7",
		"888888888888888888888888888888888888888888888888888888888888888",
		"9",
		"10--",
		"11",
		"12",
		"13",
	}

	for _, textline := range text {
		editor.AddLine()

		for _, rune := range textline {
			editor.PutRuneAtCursors(rune)
		}
	}
	editor.SetCursor(8, 0)

	ui.Loop()

}

func main() {
	start()
}
