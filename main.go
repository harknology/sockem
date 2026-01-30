package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"farthergate.com/sockem/config"
	"farthergate.com/sockem/data"
	"farthergate.com/sockem/listener"
	"farthergate.com/sockem/sender"
	"farthergate.com/sockem/state"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.BUFFER_SIZE,
	WriteBufferSize: config.BUFFER_SIZE,
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade", "error", err)
		return
	}

	slog.Info("connection established")

	channel := make(chan data.PassedMessage)

	element := state.Channels.PushBack(channel)

	go sender.SendLoop(conn, element)
	go listener.RecvLoop(conn)
}

func main() {
	if config.LOG_FORMAT == "json" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	http.HandleFunc("/ws", ws)

	slog.Info("listen", "addr", config.ListenAddr())
	log.Fatal(http.ListenAndServe(config.ListenAddr(), nil))
}
