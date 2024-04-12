package model

import "github.com/gorilla/websocket"

type WSMessage struct {
	data []byte
}

type ErrMessage struct {
	err error
}

type NewClientMessage struct {
	conn *websocket.Conn
}
