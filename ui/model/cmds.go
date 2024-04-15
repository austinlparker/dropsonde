package model

import tea "github.com/charmbracelet/bubbletea"

func showRawData(data string) tea.Cmd {
	return func() tea.Msg {
		return RawDataViewMessage{data}
	}
}
