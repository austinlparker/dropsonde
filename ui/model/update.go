package model

import (
	"encoding/json"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"slices"
)

const useHighPerformanceRenderer = false

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.msg = ""
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			Close(m.wsConn)
			return m, tea.Quit
		case "m":
			m.lastView = m.activeView
			m.activeView = ListView
			return m, nil
		case "t":
			m.activeTab = "traces"
			return m, nil
		case "l":
			m.activeTab = "logs"
			return m, nil
		case "?":
			m.lastView = m.activeView
			m.activeView = HelpView
			m.vpFullScreen = true
			return m, nil
		}
	case tea.WindowSizeMsg:
		_, h := stackListStyle.GetFrameSize()
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.tapMessageList.SetHeight(msg.Height - h - 3)

		if !m.metadataVPReady {
			m.metadataVP = viewport.New(120, m.terminalHeight/2-8)
			m.metadataVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.metadataVPReady = true
		} else {
			m.metadataVP.Width = 120
			m.metadataVP.Height = 12
		}

		if !m.valueVPReady {
			m.valueVP = viewport.New(120, m.terminalHeight-8)
			m.valueVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.valueVPReady = true
		} else {
			m.valueVP.Width = 120
			m.valueVP.Height = 12
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(120, m.terminalHeight-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(HelpText)
			m.helpVPReady = true
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
					isPresent := false
					for _, item := range m.tapMessageList.Items() {
						if kti, ok := item.(KeyTimeseriesItem); ok {
							if kti.name == ts.Name && kti.description == ts.Description {
								isPresent = true
								break
							}
						}
					}
					if !isPresent {
						m.tapMessageList.InsertItem(len(m.tapMessageList.Items()),
							KeyTimeseriesItem{
								ts:          ts,
								name:        ts.Name,
								description: ts.Description,
							})
					}
				}
			}
		}
		return m, WaitForMessages(m.channel)
	case TimeseriesChosenMessage:
		idx := slices.IndexFunc(m.metrics, func(ts Timeseries) bool {
			return ts.Name == msg.name
		})
		m.valueVP.SetContent(TimeseriesToString(m.metrics[idx]))
		return m, nil

	}

	switch m.activeView {
	case ListView:
		m.tapMessageList, cmd = m.tapMessageList.Update(msg)
		cmds = append(cmds, cmd)
	case MetadataView:
		m.metadataVP, cmd = m.metadataVP.Update(msg)
		cmds = append(cmds, cmd)
	case ValueView:
		m.valueVP, cmd = m.valueVP.Update(msg)
		cmds = append(cmds, cmd)
	case HelpView:
		m.helpVP, cmd = m.helpVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
