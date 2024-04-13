package model

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
)

type responseMsg struct {
	ResourceSpans   json.RawMessage `json:"resourceSpans,omitempty"`
	ResourceMetrics json.RawMessage `json:"resourceMetrics,omitempty"`
	ResourceLogs    json.RawMessage `json:"resourceLogs,omitempty"`
}

type delegateKeyMap struct {
	choose key.Binding
}

type KeyTimeseriesItem struct {
	ts          Timeseries
	name        string
	description string
}

func (item KeyTimeseriesItem) Title() string {
	return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim("name", 10), item.name), listWidth)
}

func (item KeyTimeseriesItem) Description() string {
	if item.description != "" {
		return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim("description", 10), item.description), listWidth)
	}
	return ""
}

func (item KeyTimeseriesItem) FilterValue() string {
	return item.name
}
