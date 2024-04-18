package components

import (
	"github.com/austinlparker/dropsonde/ui/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabItem struct {
	Name string
	Item interface{}
}

type Tabs struct {
	tabList     []TabItem
	selectedTab int
	help        HelpModel
}

func NewTabs(tabList []TabItem) *Tabs {
	return &Tabs{
		selectedTab: 0,
		help:        NewHelpModel(),
		tabList:     tabList,
	}
}

func (t *Tabs) Init() tea.Cmd { return nil }

func (t *Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.help.keys.Up), key.Matches(msg, t.help.keys.ShiftTab):
			t.decrementSelection()
		case key.Matches(msg, t.help.keys.Down), key.Matches(msg, t.help.keys.Tab):
			t.incrementSelection()
		}
	}
	return t, nil
}

func (t *Tabs) View() string {
	renderedTabs := make([]string, 0)

	for i, tl := range t.tabList {
		if i == t.selectedTab {
			renderedTabs = append(renderedTabs, theme.ActiveTabStyle.Render(tl.name))
		} else {
			renderedTabs = append(renderedTabs, theme.InactiveTabStyle.Render(tl.name))
		}
	}

	return theme.TabGroupStyle.Render(lipgloss.JoinVertical(lipgloss.Right, renderedTabs...))
}

func (t *Tabs) CurrentTab() TabItem {
	return t.tabList[t.selectedTab]
}

func (t *Tabs) decrementSelection() {
	if t.selectedTab > 0 {
		t.selectedTab--
	} else {
		t.selectedTab = len(t.tabList) - 1
	}
}

func (t *Tabs) incrementSelection() {
	if t.selectedTab == len(t.tabList)-1 {
		t.selectedTab = 0
	} else {
		t.selectedTab++
	}
}
