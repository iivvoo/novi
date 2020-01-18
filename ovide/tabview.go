package ovide

import "github.com/rivo/tview"

type Tab struct {
	Label string
	Item  tview.Primitive
}
type TabbedLayout struct {
	*tview.Flex
	Tabs   []*Tab
	Active tview.Primitive

	buttonFlex *tview.Flex
}

func NewTabbedLayout() *TabbedLayout {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	buttonFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(buttonFlex, 1, 1, false)

	flex.AddItem(tview.NewBox().SetBorder(true), 0, 0, true)
	return &TabbedLayout{
		Flex:       flex,
		buttonFlex: buttonFlex,
	}
}

func (t *TabbedLayout) AddTab(Label string, Item tview.Primitive) tview.Primitive {
	// tabs could also be a textview, similar to presentation demo
	button := tview.NewButton(Label)
	t.buttonFlex.AddItem(button, 0, 1, false)

	t.Tabs = append(t.Tabs, &Tab{Label: Label, Item: Item})

	if t.Active != nil {
		t.RemoveItem(t.Active)
	}
	t.AddItem(Item, 0, 10, true)
	t.Active = Item
	return Item
}
