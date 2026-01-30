package main

import (
	"log"
	"log/slog"
	"net/http"

	"farthergate.com/sockem/config"
	"farthergate.com/sockem/data"
	"farthergate.com/sockem/listener"
	"farthergate.com/sockem/sender"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.BUFFER_SIZE,
	WriteBufferSize: config.BUFFER_SIZE,
}

var channels = make([]chan data.PassedMessage, 0)

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade", "error", err)
		return
	}

	channel := make(chan data.PassedMessage)
	channels = append(channels, channel)

	go sender.SendLoop(conn, channel)
	go listener.RecvLoop(conn, channels)
}

func main() {
	http.HandleFunc("/ws", ws)

	slog.Info("listen", "addr", config.ListenAddr())
	log.Fatal(http.ListenAndServe(config.ListenAddr(), nil))
}
