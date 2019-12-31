package viemu

import (
	"gitlab.com/iivvoo/ovim/logger"
	"gitlab.com/iivvoo/ovim/ovim"
)

var log = logger.GetLogger("viemu")

type ViMode int

const (
	ModeEdit ViMode = iota
	ModeCommand
)

type Vi struct {
	Editor *ovim.Editor
	Mode   ViMode
}

func NewVi(e *ovim.Editor) *Vi {
	return &Vi{Editor: e, Mode: ModeCommand}
}

func (em *Vi) HandleEvent(event ovim.Event) bool {
	return false
}

func (em *Vi) GetStatus(width int) string {
	return "--INSERT-- 10,20 (fake)"
}
