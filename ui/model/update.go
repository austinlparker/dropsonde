package model

import (
	"encoding/json"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"time"
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
		case "l":
			m.lastView = m.activeView
			m.activeView = ListView
			m.vpFullScreen = false
			return m, nil
		case "?":
			m.lastView = m.activeView
			m.activeView = HelpView
			m.vpFullScreen = true
			return m, nil
		case "v":
			m.lastView = m.activeView
			m.activeView = ValueView
			m.vpFullScreen = false
			return m, nil
		case "o":
			m.lastView = m.activeView
			m.activeView = OpAmpView
			m.vpFullScreen = true
			return m, nil
		}
	case tea.WindowSizeMsg:
		h := stackListStyle.GetMaxHeight()
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.rawDataList.SetHeight(msg.Height - h - 10)

		if !m.valueVPReady {
			m.valueVP = viewport.New(120, m.terminalHeight-8)
			m.valueVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.valueVPReady = true
		} else {
			m.valueVP.Width = 120
			m.valueVP.Height = 20
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(120, m.terminalHeight-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(HelpText)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = 120
			m.helpVP.Height = 20
		}

		if !m.opampVPReady {
			m.opampVP = viewport.New(120, m.terminalHeight-7)
			m.opampVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.opampVPReady = true
		} else {
			m.opampVP.Width = 120
			m.opampVP.Height = 20
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
			var it rawDataType
			if res.ResourceSpans != nil {
				it = TraceData
			}
			if res.ResourceLogs != nil {
				it = LogData
			}
			if res.ResourceMetrics != nil {
				it = MetricData
			}
			m.rawDataList.InsertItem(len(m.rawDataList.Items()),
				RawDataItem{
					item:     string(msg.data),
					itemType: it,
					time:     time.Now(),
				})
			if res.ResourceMetrics != nil {
				metrics, err := ParseMetricMessage(msg)
				if err != nil {
					m.errorMsg = err.Error()
					return m, nil
				}
				m.metrics = AddOrUpdateMetricsToTimeseries(&metrics, m.metrics)
			}
		}
		return m, WaitForMessages(m.channel)
	case RawDataViewMessage:
		m.valueVP.SetContent(msg.data)
		return m, nil
	case OpAmpViewMessage:
		m.opampVP.SetContent(msg.data)
	}

	switch m.activeView {
	case ListView:
		m.rawDataList, cmd = m.rawDataList.Update(msg)
		cmds = append(cmds, cmd)
	case ValueView:
		m.valueVP, cmd = m.valueVP.Update(msg)
		cmds = append(cmds, cmd)
	case HelpView:
		m.helpVP, cmd = m.helpVP.Update(msg)
		cmds = append(cmds, cmd)
	case OpAmpView:
		m.opampVP, cmd = m.opampVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
