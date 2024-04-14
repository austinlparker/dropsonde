package model

import (
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type WSMessage struct {
	data []byte
}

type ErrMessage struct {
	err error
}

type ErrConnNotEstablished struct {
	err error
}

type NewClientMessage struct {
	conn *websocket.Conn
}

type CloseClientMessage struct{}

type ParsedMetricMessage struct {
	metrics pmetric.Metrics
}

type ParsedTraceMessage struct {
	traces ptrace.Traces
}

type ParsedLogMessage struct {
	logs plog.Logs
}

type TimeseriesMessage struct {
	ts []Timeseries
}

type TimeseriesChosenMessage struct {
	name string
}

type RawDataViewMessage struct {
	data string
}
