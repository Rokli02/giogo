package models

import "nhooyr.io/websocket"

type ClientData struct {
	Username string
	Conn     *websocket.Conn
}

type ClientMessage struct {
	Connection *websocket.Conn
	SocketData *SocketData
}
