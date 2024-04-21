package model

import (
	"github.com/austinlparker/dropsonde/ui/components"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
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
		case key.Matches(msg, components.KeyBindings.Help):
			m.helpView.View()
		case key.Matches(msg, components.KeyBindings.Quit):
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
		return m, tea.Batch(WaitForMessages(m.rtChannel), PushRawMessage(msg.data))
	case TableDataViewMessage:
		m.dataTable.SetContent(msg.data)
	case RawDataViewMessage:
		m.dataPager.SetContent(msg.data)
	case OpAmpViewMessage:
		m.opAmpPager.SetContent(msg.data)
	case RawDataUpdated:
		m.rawDataSlice = append(m.rawDataSlice, msg.data)
		td := RawTableData{
			columns: table.Row{"Timestamp", "Data"},
		}
		for _, item := range m.rawDataSlice {
			if item.ResourceMetrics != nil {
				td.rows = append(td.rows, table.Row{item.ts.Format("2006-01-02 15:04:05"), Trim(string(item.ResourceMetrics), 100)})
			}
			if item.ResourceLogs != nil {
				td.rows = append(td.rows, table.Row{item.ts.Format("2006-01-02 15:04:05"), Trim(string(item.ResourceLogs), 100)})
			}
			if item.ResourceSpans != nil {
				td.rows = append(td.rows, table.Row{item.ts.Format("2006-01-02 15:04:05"), Trim(string(item.ResourceSpans), 100)})
			}
		}
		m.dataTable.SetContent(td)
	default:
		m.updateAllModels(msg)
	}
	return m, cmd
}
