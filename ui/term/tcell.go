package termui

import (
	"github.com/gdamore/tcell"
	"github.com/iivvoo/ovim/ovim"
)

var KeyMap = map[tcell.Key]ovim.KeyType{
	// KeyBackspace2 is the 0x7F variant that's a regular character but > ' '
	tcell.KeyBackspace2: ovim.KeyBackspace,
	tcell.KeyEsc:        ovim.KeyEscape,
	tcell.KeyEnter:      ovim.KeyEnter,
	tcell.KeyUp:         ovim.KeyUp,
	tcell.KeyDown:       ovim.KeyDown,
	tcell.KeyLeft:       ovim.KeyLeft,
	tcell.KeyRight:      ovim.KeyRight,
	tcell.KeyHome:       ovim.KeyHome,
	tcell.KeyEnd:        ovim.KeyEnd,
	tcell.KeyPgUp:       ovim.KeyPgUp,
	tcell.KeyPgDn:       ovim.KeyPgDn,
	tcell.KeyBackspace:  ovim.KeyBackspace,
	tcell.KeyTab:        ovim.KeyTab,
	tcell.KeyDelete:     ovim.KeyDelete,
	tcell.KeyInsert:     ovim.KeyInsert,
	tcell.KeyF1:         ovim.KeyF1,
	tcell.KeyF2:         ovim.KeyF2,
	tcell.KeyF3:         ovim.KeyF3,
	tcell.KeyF4:         ovim.KeyF4,
	tcell.KeyF5:         ovim.KeyF5,
	tcell.KeyF6:         ovim.KeyF6,
	tcell.KeyF7:         ovim.KeyF7,
	tcell.KeyF8:         ovim.KeyF8,
	tcell.KeyF9:         ovim.KeyF9,
	tcell.KeyF10:        ovim.KeyF10,
	tcell.KeyF11:        ovim.KeyF11,
	tcell.KeyF12:        ovim.KeyF12,
}

type DecomposedKey struct {
	Modifier ovim.KeyModifier
	Key      ovim.KeyType
	Rune     rune
}

var DecomposeMap = map[tcell.Key]DecomposedKey{
	tcell.KeyCtrlSpace: DecomposedKey{ovim.ModCtrl, ovim.KeyRune, ' '},
	tcell.KeyCtrlA:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'a'},
	tcell.KeyCtrlB:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'b'},
	tcell.KeyCtrlC:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'c'},
	tcell.KeyCtrlD:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'd'},
	tcell.KeyCtrlE:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'e'},
	tcell.KeyCtrlF:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'f'},
	tcell.KeyCtrlG:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'g'},
	tcell.KeyCtrlH:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'h'},
	tcell.KeyCtrlI:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'i'},
	tcell.KeyCtrlJ:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'j'},
	tcell.KeyCtrlK:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'k'},
	tcell.KeyCtrlL:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'l'},
	tcell.KeyCtrlM:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'm'},
	tcell.KeyCtrlN:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'n'},
	tcell.KeyCtrlO:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'o'},
	tcell.KeyCtrlP:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'p'},
	tcell.KeyCtrlQ:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'q'},
	tcell.KeyCtrlR:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'r'},
	tcell.KeyCtrlS:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 's'},
	tcell.KeyCtrlT:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 't'},
	tcell.KeyCtrlU:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'u'},
	tcell.KeyCtrlV:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'v'},
	tcell.KeyCtrlW:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'w'},
	tcell.KeyCtrlX:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'x'},
	tcell.KeyCtrlY:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'y'},
	tcell.KeyCtrlZ:     DecomposedKey{ovim.ModCtrl, ovim.KeyRune, 'z'},
	// ctrl-Escape, according to tcell?
	tcell.KeyCtrlLeftSq:    DecomposedKey{ovim.ModCtrl, ovim.KeyEscape, ' '},
	tcell.KeyCtrlBackslash: DecomposedKey{ovim.ModCtrl, ovim.KeyRune, '\\'},
	// tcell.KeyCtrlRightSq:    DecomposedKey{ovim.ModCtrl, ovim.KeyRune, ' '},
	tcell.KeyCtrlCarat:      DecomposedKey{ovim.ModCtrl, ovim.KeyRune, '^'},
	tcell.KeyCtrlUnderscore: DecomposedKey{ovim.ModCtrl, ovim.KeyRune, '_'},
}

func MapTCellKey(ev *tcell.EventKey) ovim.Event {
	/*
	   type EventKey struct {
	   	t   time.Time
	   	mod ModMask
	   	key Key
	   	ch  rune
	   }

	   rune = character, mod=alt, shift, etc
	   Key() is smashed, eg CtrlL

	   constants:
	   https://github.com/gdamore/tcell/blob/de7e78efa4a71b3f36c7154989c529dbdf9ae623/key.go#L83

	   use map to map types
	   Terminals don't (always) decompose modifiers. For example,
	   'H' is just H, not shift-h
	   but also a lot of contol characters are sent as-is,
	   so CtrlL is not ctrl-l

	   We can decompose this and transform CtrlL into ctrl-l

	*/
	key := ev.Key()
	if ovimKey, ok := KeyMap[key]; ok {
		return &ovim.KeyEvent{Key: ovimKey}
	} else if decomposed, ok := DecomposeMap[key]; ok {
		return &ovim.KeyEvent{Modifier: decomposed.Modifier, Key: decomposed.Key, Rune: decomposed.Rune}
	} else {
		return &ovim.CharacterEvent{Rune: ev.Rune()}
	}
}
