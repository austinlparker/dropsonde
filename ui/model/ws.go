package model

import (
	"encoding/json"
	"github.com/austinlparker/dropsonde/internal/parsers"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ResponseMsg struct {
	ResourceSpans   json.RawMessage `json:"resourceSpans,omitempty"`
	ResourceMetrics json.RawMessage `json:"resourceMetrics,omitempty"`
	ResourceLogs    json.RawMessage `json:"resourceLogs,omitempty"`
	ts              time.Time
}

type RawTableData struct {
	columns table.Row
	rows    []table.Row
}

func (rd RawTableData) Columns() table.Row { return rd.columns }
func (rd RawTableData) Rows() []table.Row  { return rd.rows }

func Dial(endpoint string) tea.Cmd {
	return func() tea.Msg {
		conn, err := newConn(endpoint)
		if err != nil {
			return ErrMessage{err}
		}
		return NewClientMessage{conn}
	}
}

func Close(ws *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		closeConn(ws)
		return CloseClientMessage{}
	}
}

func ListenForMessage(conn *websocket.Conn, c chan []byte) tea.Cmd {
	return func() tea.Msg {
		for {
			if conn == nil {
				return ErrConnNotEstablished{nil}
			}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return ErrConnNotEstablished{err}
				//TODO: handle reconnection
			}
			c <- msg
		}
	}
}

func WaitForMessages(c chan []byte) tea.Cmd {
	return func() tea.Msg {
		return WSMessage{<-c}
	}
}

func newConn(endpoint string) (*websocket.Conn, error) {
	wsUrl := url.URL{Scheme: "ws", Host: endpoint, Path: "/"}
	headers := make(http.Header)
	headers.Set("Origin", "http://localhost")
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), headers)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeConn(conn *websocket.Conn) {
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("error:", err)
		return
	}
}

func PushRawMessage(msg []byte) tea.Cmd {
	return func() tea.Msg {
		tea.Println("Raw Message:", string(msg))
		var res ResponseMsg
		err := json.Unmarshal(msg, &res)
		if err != nil {
			return ErrMessage{err}
		}
		tea.Println("pushing raw message")
		res.ts = time.Now()
		return RawDataUpdated{res}
	}
}

func ParseMessage(msg []byte) tea.Cmd {
	return func() tea.Msg {
		var res ResponseMsg
		err := json.Unmarshal(msg, &res)
		if err != nil {
			return ErrMessage{err}
		}
		if res.ResourceMetrics != nil {
			metric, err := parsers.ParseMetricMessage(msg)
			if err != nil {
				return ErrMessage{err}
			}
			return ParsedMetricMessage{metric}
		}
		if res.ResourceSpans != nil {
			spans, err := parsers.ParseTraceMessage(msg)
			if err != nil {
				return ErrMessage{err}
			}
			return ParsedTraceMessage{spans}
		}
		if res.ResourceLogs != nil {
			return ErrMessage{nil}
		}
		return ErrMessage{nil}
	}
}
