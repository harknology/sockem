package sender

import (
	"encoding/json"
	"log/slog"
	"slices"

	"farthergate.com/sockem/data"
	"github.com/gorilla/websocket"
)

func SendLoop(c *websocket.Conn, messages chan data.PassedMessage) {
	defer c.Close()

	var listeners []string

	for {
		msg := <-messages

		if msg.InternalMessage != data.NONE {
			switch msg.InternalMessage {
			case data.LISTEN:
				if !msg.Authenticated {
					c.WriteMessage(websocket.TextMessage, []byte(`{"error":"not authenticated"}`))
					continue
				}
				listeners = append(listeners, msg.Parsed.Data.(string))
			}
		} else {
			if !slices.Contains(listeners, msg.Parsed.Name) {
				continue
			}

			serialized, err := json.Marshal(msg.Parsed)
			if err != nil {
				slog.Error("serialize", "error", err)
				continue
			}

			c.WriteMessage(websocket.TextMessage, serialized)
		}
	}
}
