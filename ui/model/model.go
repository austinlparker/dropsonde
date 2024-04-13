package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type stateView uint

const (
	ListView stateView = iota
	MetadataView
	ValueView
	HelpView
)

type model struct {
	tapEndpoint       string
	opAmpEndpoint     string
	tabs              []string
	activeView        stateView
	tapMessageList    list.Model
	activeTab         string
	wsConn            *websocket.Conn
	msg               string
	errorMsg          string
	traces            []ptrace.Traces
	metrics           []Timeseries
	logs              []plog.Logs
	shouldQuit        bool
	channel           chan []byte
	valueVP           viewport.Model
	metadataVP        viewport.Model
	helpVP            viewport.Model
	valueVPReady      bool
	metadataVPReady   bool
	helpVPReady       bool
	vpFullScreen      bool
	showHelpIndicator bool
	debugMode         bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		Dial(m.tapEndpoint),
		WaitForMessages(m.channel),
	)
}
