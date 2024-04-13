package model

import (
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"slices"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	m.msg = ""
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "m":
			m.activeTab = "metrics"
			return m, nil
		case "t":
			m.activeTab = "traces"
			return m, nil
		case "l":
			m.activeTab = "logs"
			return m, nil
		}
	case NewClientMessage:
		m.wsConn = msg.conn
		return m, ListenForMessage(m.wsConn, m.channel)
	case ErrConnNotEstablished:
		Dial(m.tapEndpoint)
		return m, ListenForMessage(m.wsConn, m.channel)
	case WSMessage:
		if msg.data != nil {
			var res responseMsg
			err := json.Unmarshal(msg.data, &res)
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			if res.ResourceMetrics != nil {
				metrics, err := ParseMetricMessage(msg)
				if err != nil {
					m.errorMsg = err.Error()
					return m, nil
				}
				m.metrics = AddOrUpdateMetricsToTimeseries(&metrics, m.metrics)
				for _, ts := range m.metrics {
					m.tapMessageList.InsertItem(len(m.tapMessageList.Items()),
						KeyTimeseriesItem{
							ts:          ts,
							name:        ts.Name,
							description: ts.Description,
						})
				}
			}
		}
		return m, WaitForMessages(m.channel)
	case TimeseriesChosenMessage:
		idx := slices.IndexFunc(m.metrics, func(ts Timeseries) bool {
			return ts.Name == msg.name
		})
		content, err := json.MarshalIndent(m.metrics[idx], "", "    ")
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.valueVP.SetContent(string(content))
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
