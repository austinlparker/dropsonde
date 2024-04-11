package main

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type responseMsg struct {
	ResourceSpans   json.RawMessage `json:"resourceSpans,omitempty"`
	ResourceMetrics json.RawMessage `json:"resourceMetrics,omitempty"`
	ResourceLogs    json.RawMessage `json:"resourceLogs,omitempty"`
}

type tapMessage struct {
	data []byte
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func NewWsConnection() (*websocket.Conn, error) {
	wsURL := url.URL{Scheme: "ws", Host: "localhost:12001", Path: "/"}
	headers := make(http.Header)
	headers.Set("Origin", "http://localhost")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), headers)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func waitForMessages(ws chan []byte) tea.Cmd {
	return func() tea.Msg {
		return tapMessage{
			<-ws,
		}
	}
}

func listenForMessages(conn *websocket.Conn, c chan []byte) tea.Cmd {
	return func() tea.Msg {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				tea.Println("Error reading message", err)
			}
			c <- msg
		}
	}
}

func (m *model) assignMessage(msg tapMessage) {
	if msg.data != nil {
		var response responseMsg
		err := json.Unmarshal(msg.data, &response)
		if err != nil {
			tea.Println("Error unmarshaling message", err)
			return
		}
		if response.ResourceSpans != nil {
			unmarshaler := &ptrace.JSONUnmarshaler{}
			span, err := unmarshaler.UnmarshalTraces(msg.data)
			if err != nil {
				tea.Println("Error unmarshaling message", err)
				return
			}
			m.traces = append(m.traces, span)
		}
		if response.ResourceMetrics != nil {
			unmarshaler := &pmetric.JSONUnmarshaler{}
			metrics, err := unmarshaler.UnmarshalMetrics(msg.data)
			if err != nil {
				tea.Println("Error unmarshaling message", err)
				return
			}
			m.metrics = append(m.metrics, metrics)
		}
		if response.ResourceLogs != nil {
			unmarshaler := &plog.JSONUnmarshaler{}
			log, err := unmarshaler.UnmarshalLogs(msg.data)
			if err != nil {
				tea.Println("Error unmarshaling message", err)
				return
			}
			m.logs = append(m.logs, log)
		}
	}
}

type model struct {
	tabs       []string
	activeTab  int
	wsConn     *websocket.Conn
	msg        string
	traces     []ptrace.Traces
	metrics    []pmetric.Metrics
	logs       []plog.Logs
	shouldQuit bool
	channel    chan []byte
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForMessages(m.wsConn, m.channel),
		waitForMessages(m.channel),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := tmsg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	case tapMessage:
		m.assignMessage(msg.(tapMessage))
		return m, waitForMessages(m.channel)
	}
	return m, nil
}

func (m model) View() string {
	view := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	row := lipgloss.NewStyle().Width(80).Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	view.WriteString(row)
	view.WriteString("\n")
	view.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.renderMetrics()))
	return docStyle.Render(view.String())
}

func (m model) renderMetrics() string {
	metricString := strings.Builder{}
	for i := 0; i < len(m.metrics); i++ {
		metric := m.metrics[i]
		for j := 0; j < metric.ResourceMetrics().Len(); j++ {
			resourceMetric := metric.ResourceMetrics().At(j)
			scopeMetrics := resourceMetric.ScopeMetrics()
			for k := 0; k < scopeMetrics.Len(); k++ {
				scopeMetric := scopeMetrics.At(k)
				name := scopeMetric.Scope().Name()
				metricString.WriteString(name)
				for v := 0; v < scopeMetrics.Len(); v++ {
					metrics := scopeMetric.Metrics()
					for l := 0; l < metrics.Len(); l++ {
						value := metrics.At(l)
						metricString.WriteString(value.Name())
						metricString.WriteString("\n")
					}
				}
			}
		}
	}
	return metricString.String()
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	wsConn, err := NewWsConnection()
	if err != nil {
		fmt.Println("could not connect to websocket:", err)
		os.Exit(1)
	}

	tabs := []string{"Metrics", "Traces", "Logs"}

	p := tea.NewProgram(&model{
		wsConn:  wsConn,
		channel: make(chan []byte),
		tabs:    tabs,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
