package model

import "github.com/charmbracelet/bubbles/list"

func Initial(tapEndpoint string, opAmpEndpoint string) model {
	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)
	listItems := make([]list.Item, 0)

	tabs := []string{"Metrics", "Traces", "Logs"}
	m := model{
		tapEndpoint:    tapEndpoint,
		opAmpEndpoint:  opAmpEndpoint,
		tapMessageList: list.New(listItems, appDelegate, listWidth+10, 0),
		tabs:           tabs,
		channel:        make(chan []byte),
	}

	m.tapMessageList.Title = "OpenTelemetry Metrics"
	m.tapMessageList.SetStatusBarItemName("metric", "metrics")
	m.tapMessageList.SetFilteringEnabled(false)
	m.tapMessageList.SetShowHelp(false)

	return m
}
