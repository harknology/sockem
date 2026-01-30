package data

import "github.com/gorilla/websocket"

type Message struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

const (
	NONE         uint8 = 0
	AUTHENTICATE uint8 = 1
	LISTEN       uint8 = 2
	CLOSE        uint8 = 3
)

type PassedMessage struct {
	InternalMessage uint8
	Authenticated   bool
	Parsed          Message
	Conn            *websocket.Conn
}
