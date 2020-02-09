package ovide

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

// is it a TabView, TabLayout, ??
// Flex, Frame and Grid are also just that. "Tabs"?
//
type Tab struct {
	Id    string
	Label string
	Item  tview.Primitive
}

type Tabs struct {
	*tview.Flex
	Tabs   []*Tab
	Active *Tab

	labels *tview.TextView
}

func NewTabs() *Tabs {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	labels := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)
	flex.AddItem(labels, 1, 1, false)

	flex.AddItem(tview.NewBox().SetBorder(true), 0, 0, true)
	return &Tabs{
		Flex:   flex,
		labels: labels,
	}
}

func (t *Tabs) updateLabels() {
	s := ""
	for i, tab := range t.Tabs {
		s += fmt.Sprintf(`["%d"][darkcyan]%s[white][""]  `, i, tab.Label)
		if tab.Id == t.Active.Id {
			t.labels.Highlight(strconv.Itoa(i))
		}
	}
	t.labels.SetText(s)
}

func (t *Tabs) AddTab(Id string, Label string, Item tview.Primitive) tview.Primitive {
	// If a tab with this Id already exists, select it in stead
	// (but it's better if the caller would do that check)
	if t.SelectTab(Id) {
		return nil
	}

	tab := &Tab{Id: Id, Label: Label, Item: Item}
	t.Tabs = append(t.Tabs, tab)
	t.setActive(tab)
	return Item
}

func (t *Tabs) setActive(tab *Tab) {
	if t.Active != nil {
		t.RemoveItem(t.Active.Item)
	}
	t.AddItem(tab.Item, 0, 10, true)
	t.Active = tab
	t.updateLabels()
}

func (t *Tabs) SelectTab(Id string) bool {
	for _, tab := range t.Tabs {
		if tab.Id == Id {
			t.setActive(tab)
			return true
		}
	}
	return false
}

func (t *Tabs) CloseTab(Id string) bool {
	for i, tab := range t.Tabs {
		if tab.Id == Id {
			t.Tabs = append(t.Tabs[:i], t.Tabs[i+1:]...)
			t.RemoveItem(tab.Item)
			t.updateLabels()
			return true
		}
	}
	return false

}
