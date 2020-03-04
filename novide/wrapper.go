package novide

import (
	"time"

	"github.com/iivvoo/novi/novi"
	"github.com/rivo/tview"
)

// OviWrapper wraps the OviPrimitive into something that we can pass to the novi.Core
type OviWrapper struct {
	prim *Ovi
	app  *tview.Application
}

// NewWrapper wraps an Ovi primitive
func NewWrapper(app *tview.Application, prim *Ovi) *OviWrapper {
	return &OviWrapper{prim: prim, app: app}
}

// Finish is called when the editor finishes
func (o *OviWrapper) Finish() {}

// Loop allows the UI to start a loop goroutine
func (o *OviWrapper) Loop(c chan novi.Event) {
	o.prim.SetChan(c)
}

// Render is called to render the UI if necessary
func (o *OviWrapper) Render() {}

// GetDimension returns the dimension of the editor
func (o *OviWrapper) GetDimension() (int, int) {
	return o.prim.GetDimension()
}

// AskInput will instruct the UI to ask for additional, "inline" input
func (o *OviWrapper) AskInput(string) novi.InputSource {
	// handle keys from status
	o.prim.Source = CommandSource

	o.UpdateInput(CommandSource, "", 0)
	return CommandSource
}

// CloseInput closes the inline input
func (o *OviWrapper) CloseInput(novi.InputSource) {
	o.prim.Source = MainSource
}

// UpdateInput is called to update the inline input
func (o *OviWrapper) UpdateInput(source novi.InputSource, s string, pos int) {
	o.prim.UpdateInput(":"+s, pos+1)
}

// SetStatus sets the status of the editor
func (o *OviWrapper) SetStatus(status string) {
	o.prim.UpdateStatus(status)
}

// SetError sets the error on the input (and should clear it after a while)
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
