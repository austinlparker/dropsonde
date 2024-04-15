package model

import (
	"github.com/charmbracelet/bubbles/list"
	"os"
)

func Initial(tapEndpoint string) model {
	var rawDelegateKeys = newRawDelegateKeyMap()
	rawDelegate := newRawItemDelegate(rawDelegateKeys)
	rawListItems := make([]list.Item, 0)

	var dbg bool
	if len(os.Getenv("DEBUG")) > 0 {
		dbg = true
	}
	var opampsrv Server
	NewServer(&opampsrv)
	m := model{
		tapEndpoint:       tapEndpoint,
		rawDataList:       list.New(rawListItems, rawDelegate, listWidth+10, 0),
		channel:           make(chan []byte),
		debugMode:         dbg,
		showHelpIndicator: true,
		useRawDataView:    true,
		opamp:             opampsrv,
	}

	m.rawDataList.Title = "Remote Tap"
	m.rawDataList.SetStatusBarItemName("item", "items")
	m.rawDataList.SetFilteringEnabled(true)
	m.rawDataList.SetShowHelp(false)
	m.rawDataList.SetShowPagination(false)

	return m
}
