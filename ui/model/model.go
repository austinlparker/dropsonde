package model

import (
	"github.com/austinlparker/dropsonde/ui/components"
	"github.com/austinlparker/dropsonde/ui/opamp"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

type model struct {
	tabs *components.Tabs

	dataPager  components.ContentModel
	opAmpPager components.ContentModel
	dataTable  components.ContentModel
	helpView   components.HelpModel

	isQuitting   bool
	errorMessage string

	spinner spinner.Model

	width, height int

	rtEndpoint string
	rtChannel  chan []byte
	wsConn     *websocket.Conn
	oaServer   opamp.Server

	rawDataSlice []ResponseMsg
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		Dial(m.rtEndpoint),
		WaitForMessages(m.rtChannel),
	)
}

func NewModel(tapEndpoint string) *model {
	var opampsrv opamp.Server
	opamp.NewServer(&opampsrv)
	m := &model{
		rtEndpoint:   tapEndpoint,
		rtChannel:    make(chan []byte),
		oaServer:     opampsrv,
		dataPager:    components.NewPager(),
		opAmpPager:   components.NewPager(),
		dataTable:    components.NewTable([]int{10, 90}),
		helpView:     components.NewHelpModel(),
		isQuitting:   false,
		width:        0,
		height:       0,
		spinner:      components.NewSpinner(),
		rawDataSlice: make([]ResponseMsg, 0),
	}
	m.tabs = components.NewTabs([]components.TabItem{
		{Name: "raw", Item: m.dataTable},
		{Name: "opamp", Item: m.opAmpPager},
	})
	return m
}
