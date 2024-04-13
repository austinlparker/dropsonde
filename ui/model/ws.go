package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

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
