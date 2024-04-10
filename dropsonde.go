package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type responseMsg []byte

type wsConnection struct {
	conn    *websocket.Conn
	channel chan []byte
}

func NewWsConnection() (*wsConnection, error) {
	url := url.URL{Scheme: "ws", Host: "localhost:12001", Path: "/"}
	headers := make(http.Header)
	headers.Set("Origin", "http://localhost")
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	if err != nil {
		return nil, err
	}

	wsConn := &wsConnection{
		conn:    conn,
		channel: make(chan []byte),
	}

	go func() {
		for {
			_, message, err := wsConn.conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}
			wsConn.channel <- message
		}
	}()

	return wsConn, nil
}

func waitForActivity(m *model) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-m.wsConn.channel)
	}
}

func (m model) assignMessage(msg []byte) {
	var data map[string]interface{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println("error decoding JSON:", err)
		return
	}

	for key := range data {
		switch key {
		case "resourceSpans":
			unmarshaler := &ptrace.JSONUnmarshaler{}
			span, err := unmarshaler.UnmarshalTraces(data[key].([]byte))
			if err != nil {
				fmt.Println("error decoding JSON:", err)
			}
			m.traces = append(m.traces, span)
		case "resourceMetrics":
			unmarshaler := &pmetric.JSONUnmarshaler{}
			metrics, err := unmarshaler.UnmarshalMetrics(data[key].([]byte))
			if err != nil {
				fmt.Println("error decoding JSON:", err)
			}
			metrics.
				m.metrics = append(m.metrics, metrics)
		case "resourceLogs":
			unmarshaler := &plog.JSONUnmarshaler{}
			log, err := unmarshaler.UnmarshalLogs(data[key].([]byte))
			if err != nil {
				fmt.Println("error decoding JSON:", err)
			}
			m.logs = append(m.logs, log)
		default:
			return
		}
	}
}

type model struct {
	tabs       []string
	activeTab  int
	wsConn     *wsConnection
	msg        string
	traces     []ptrace.Traces
	metrics    []pmetric.Metrics
	logs       []plog.Logs
	shouldQuit bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(&m), // wait for activity
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.shouldQuit = true
		return m, tea.Quit
	case responseMsg:
		return m, waitForActivity(&m) // wait for next event
	default:
		return m, nil
	}
}

func (m model) View() string {
	for _, metric := range m.metrics {

		s += fmt.Sprintf("Metric Name: %s\n", metric)
	}
	s += "\n Press any key to exit\n"
	if m.shouldQuit {
		s += "\n"
	}
	return s
}

func main() {
	wsConn, err := NewWsConnection()
	if err != nil {
		fmt.Println("could not connect to websocket:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(model{
		wsConn: wsConn,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
