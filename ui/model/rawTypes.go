package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"strconv"
	"time"
)

type rawDataType uint

const (
	MetricData rawDataType = iota
	TraceData
	LogData
)

type rawDelegateKeyMap struct {
	details key.Binding
}

type RawDataItem struct {
	item     string
	itemType rawDataType
	time     time.Time
}

func (item RawDataItem) Title() string {
	return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim("item", 10), item.item), listWidth)
}

func (item RawDataItem) Description() string {
	return item.time.String()
}

func (item RawDataItem) FilterValue() string {
	return strconv.Itoa(int(item.itemType))
}
