package model

import tea "github.com/charmbracelet/bubbletea"

func showItemDetails(name string) tea.Cmd {
	return func() tea.Msg {
		return TimeseriesChosenMessage{name}
	}
}
