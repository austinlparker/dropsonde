package model

import (
	"github.com/austinlparker/dropsonde/ui"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.helpView.SetWidth(msg.Width)
		m.updateAllModels(msg)
		return m, m.refresh()

	case tea.KeyMsg:
		m.tabs.Update(msg)
		switch {
		case key.Matches(msg, ui.KeyBindings.Help):
			m.helpView.View()
		case key.Matches(msg, ui.KeyBindings.Quit):
			m.isQuitting = true
			closeConn(m.wsConn)
			return m, tea.Quit
		}
		activeTab := m.ActiveTab()
		activeTab.Update(msg)
	case NewClientMessage:
		m.wsConn = msg.conn
		return m, ListenForMessage(m.wsConn, m.rtChannel)
	case ErrConnNotEstablished:
		Dial(m.rtEndpoint)
		return m, ListenForMessage(m.wsConn, m.rtChannel)
	case WSMessage:
		parseWSMessageToTable(msg.data)
		return m, WaitForMessages(m.rtChannel)
	case TableDataViewMessage:
		m.dataTable.SetContent(msg.data)
	case RawDataViewMessage:
		m.dataPager.SetContent(msg.data)
	case OpAmpViewMessage:
		m.opAmpPager.SetContent(msg.data)
	}

	return m, cmd
}
