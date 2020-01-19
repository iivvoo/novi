package ovide

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Modal struct {
	Success func(string)
	Cancel  func(string)

	app   *tview.Application
	pages *tview.Pages
	input *tview.InputField
}

// Not (new) FileModal but (NewFile) modal!
func NewFileModal(app *tview.Application, pages *tview.Pages) *Modal {
	m := &Modal{app: app, pages: pages}

	m.input = tview.NewInputField().
		SetLabel("Filename: ").
		SetFieldWidth(30).
		// some func that disallows /
		// SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(m.doneFunc)

	m.input.SetBorder(true)
	m.input.SetTitle("New file")
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}
	pages.AddPage("modal", modal(m.input, 40, 3), true, true)
	app.SetFocus(m.input)

	return m
}

func (m *Modal) doneFunc(key tcell.Key) {
	m.pages.RemovePage("modal")
	inp := m.input.GetText()
	if key == tcell.KeyEscape {
		if m.Cancel != nil {
			m.Cancel(inp)
		}
	} else if key == tcell.KeyEnter {
		if m.Success != nil {
			m.Success(inp)
		}
	}
}

func (m *Modal) SetSuccess(f func(string)) *Modal {
	m.Success = f
	return m
}

func (m *Modal) SetCancel(f func(string)) *Modal {
	m.Cancel = f
	return m
}
