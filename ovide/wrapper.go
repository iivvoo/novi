package ovide

import (
	"time"

	"github.com/iivvoo/ovim/ovim"
	"github.com/rivo/tview"
)

// OviWrapper wraps the OviPrimitive into something that we can pass to the ovim.Core
type OviWrapper struct {
	prim *Ovi
	app  *tview.Application
}

func NewWrapper(app *tview.Application, prim *Ovi) *OviWrapper {
	return &OviWrapper{prim: prim, app: app}
}

func (o *OviWrapper) Finish() {}
func (o *OviWrapper) Loop(c chan ovim.Event) {
	o.prim.SetChan(c)
}

func (o *OviWrapper) Render() {
}
func (o *OviWrapper) GetDimension() (int, int) {
	return o.prim.GetDimension()
}

func (o *OviWrapper) AskInput(string) ovim.InputSource {
	// handle keys from status
	o.prim.Source = CommandSource

	o.UpdateInput(CommandSource, "", 0)
	return CommandSource
}
func (o *OviWrapper) CloseInput(ovim.InputSource) {
	o.prim.Source = MainSource
}

func (o *OviWrapper) UpdateInput(source ovim.InputSource, s string, pos int) {
	o.prim.UpdateInput(":"+s, pos+1)
}

func (o *OviWrapper) SetStatus(status string) {
	o.prim.UpdateStatus(status)
}

func (o *OviWrapper) SetError(error string) {
	o.prim.UpdateError(error)
	go func() {
		time.Sleep(time.Second * 3)
		o.app.QueueUpdateDraw(func() {
			o.SetError("")
			// or store current status ourselves?
			o.SetStatus(o.prim.statusMsg)
			log.Printf("Error cleared")
		})
	}()
}
