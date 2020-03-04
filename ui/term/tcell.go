package termui

import (
	"github.com/gdamore/tcell"
	"github.com/iivvoo/novi/novi"
)

var KeyMap = map[tcell.Key]novi.KeyType{
	// KeyBackspace2 is the 0x7F variant that's a regular character but > ' '
	tcell.KeyBackspace2: novi.KeyBackspace,
	tcell.KeyEsc:        novi.KeyEscape,
	tcell.KeyEnter:      novi.KeyEnter,
	tcell.KeyUp:         novi.KeyUp,
	tcell.KeyDown:       novi.KeyDown,
	tcell.KeyLeft:       novi.KeyLeft,
	tcell.KeyRight:      novi.KeyRight,
	tcell.KeyHome:       novi.KeyHome,
	tcell.KeyEnd:        novi.KeyEnd,
	tcell.KeyPgUp:       novi.KeyPgUp,
	tcell.KeyPgDn:       novi.KeyPgDn,
	tcell.KeyBackspace:  novi.KeyBackspace,
	tcell.KeyTab:        novi.KeyTab,
	tcell.KeyDelete:     novi.KeyDelete,
	tcell.KeyInsert:     novi.KeyInsert,
	tcell.KeyF1:         novi.KeyF1,
	tcell.KeyF2:         novi.KeyF2,
	tcell.KeyF3:         novi.KeyF3,
	tcell.KeyF4:         novi.KeyF4,
	tcell.KeyF5:         novi.KeyF5,
	tcell.KeyF6:         novi.KeyF6,
	tcell.KeyF7:         novi.KeyF7,
	tcell.KeyF8:         novi.KeyF8,
	tcell.KeyF9:         novi.KeyF9,
	tcell.KeyF10:        novi.KeyF10,
	tcell.KeyF11:        novi.KeyF11,
	tcell.KeyF12:        novi.KeyF12,
}

type DecomposedKey struct {
	Modifier novi.KeyModifier
	Key      novi.KeyType
	Rune     rune
}

var DecomposeMap = map[tcell.Key]DecomposedKey{
	tcell.KeyCtrlSpace: DecomposedKey{novi.ModCtrl, novi.KeyRune, ' '},
	tcell.KeyCtrlA:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'a'},
	tcell.KeyCtrlB:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'b'},
	tcell.KeyCtrlC:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'c'},
	tcell.KeyCtrlD:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'd'},
	tcell.KeyCtrlE:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'e'},
	tcell.KeyCtrlF:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'f'},
	tcell.KeyCtrlG:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'g'},
	tcell.KeyCtrlH:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'h'},
	tcell.KeyCtrlI:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'i'},
	tcell.KeyCtrlJ:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'j'},
	tcell.KeyCtrlK:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'k'},
	tcell.KeyCtrlL:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'l'},
	tcell.KeyCtrlM:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'm'},
	tcell.KeyCtrlN:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'n'},
	tcell.KeyCtrlO:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'o'},
	tcell.KeyCtrlP:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'p'},
	tcell.KeyCtrlQ:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'q'},
	tcell.KeyCtrlR:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'r'},
	tcell.KeyCtrlS:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 's'},
	tcell.KeyCtrlT:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 't'},
	tcell.KeyCtrlU:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'u'},
	tcell.KeyCtrlV:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'v'},
	tcell.KeyCtrlW:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'w'},
	tcell.KeyCtrlX:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'x'},
	tcell.KeyCtrlY:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'y'},
	tcell.KeyCtrlZ:     DecomposedKey{novi.ModCtrl, novi.KeyRune, 'z'},
	// ctrl-Escape, according to tcell?
	tcell.KeyCtrlLeftSq:    DecomposedKey{novi.ModCtrl, novi.KeyEscape, ' '},
	tcell.KeyCtrlBackslash: DecomposedKey{novi.ModCtrl, novi.KeyRune, '\\'},
	// tcell.KeyCtrlRightSq:    DecomposedKey{novi.ModCtrl, novi.KeyRune, ' '},
	tcell.KeyCtrlCarat:      DecomposedKey{novi.ModCtrl, novi.KeyRune, '^'},
	tcell.KeyCtrlUnderscore: DecomposedKey{novi.ModCtrl, novi.KeyRune, '_'},
}

func MapTCellKey(ev *tcell.EventKey) novi.Event {
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
	if noviKey, ok := KeyMap[key]; ok {
		return &novi.KeyEvent{Key: noviKey}
	} else if decomposed, ok := DecomposeMap[key]; ok {
		return &novi.KeyEvent{Modifier: decomposed.Modifier, Key: decomposed.Key, Rune: decomposed.Rune}
	} else {
		return &novi.CharacterEvent{Rune: ev.Rune()}
	}
}
