package model

import tea "github.com/charmbracelet/bubbletea"

func showRawData(data string) tea.Cmd {
	return func() tea.Msg {
		return RawDataViewMessage{data}
	}
}

func (m *model) showOpAmpData() tea.Cmd {
	return func() tea.Msg {
		op := m.oaServer.GetAgents()
		return OpAmpViewMessage{op}
	}
}
