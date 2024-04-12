package model

import (
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type model struct {
	tapEndpoint   string
	opAmpEndpoint string
	tabs          []string
	activeTab     int
	wsConn        *websocket.Conn
	msg           string
	traces        []ptrace.Traces
	metrics       []pmetric.Metrics
	logs          []plog.Logs
	shouldQuit    bool
	channel       chan []byte
}
