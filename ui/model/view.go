package model

import (
	"fmt"
	"github.com/austinlparker/dropsonde/ui/components"
	"github.com/austinlparker/dropsonde/ui/theme"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() string {
	if m.isQuitting {
		return ""
	}

	var v string

	activeTab := m.ActiveTab()

	if m.errorMessage != "" {
		v = lipgloss.NewStyle().Foreground(theme.DefaultTheme().Accent).Render(m.errorMessage)
	} else if activeTab.IsReady() {
		v = activeTab.View()
	} else { // show spinner if tab's data is not ready
		v = fmt.Sprintf("\n %s \n\n", m.spinner.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		components.Header(),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Right,
				m.tabs.View(),
			),
			theme.ActiveItemStyle.Render(v),
		),
		m.helpView.View(),
	)
}

func (m *model) setUnreadyAllModels() {
	m.opAmpPager.SetUnready()
	m.dataPager.SetUnready()
	m.dataTable.SetUnready()
}

func (m *model) updateAllModels(msg tea.Msg) {
	m.opAmpPager.Update(msg)
	m.dataPager.Update(msg)
	m.dataTable.Update(msg)
}

func (m *model) refresh() tea.Cmd {
	m.errorMessage = ""

	return tea.Batch(
		components.SetModelLoading,
	)
}

func (m *model) ActiveTab() components.ContentModel {
	item := m.tabs.CurrentTab().Item
	cm := item.(components.ContentModel)
	return cm
}
