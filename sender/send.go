package sender

import (
	"encoding/json"
	"log/slog"
	"slices"
	"strings"

	"farthergate.com/sockem/data"
	"github.com/gorilla/websocket"
)

func SendLoop(c *websocket.Conn, messages chan data.PassedMessage) {
	defer c.Close()

	listeners := make([]string, 0)

	for {
		msg := <-messages
		if msg.InternalMessage == data.CLOSE && msg.Conn == c {
			break
		}

		if msg.InternalMessage != data.NONE && msg.Conn == c {
			switch msg.InternalMessage {
			case data.LISTEN:
				if !msg.Authenticated {
					c.WriteMessage(websocket.TextMessage, []byte(`{"error":"not authenticated"}`))
					continue
				}
				if strings.HasPrefix(msg.Parsed.Data.(string), "__sockem:") {
					c.WriteMessage(websocket.TextMessage, []byte(`{"error":"reserved channel"}`))
					continue
				}
				listeners = append(listeners, msg.Parsed.Data.(string))
				slog.Info("listener added", "channel", msg.Parsed.Data.(string))
				c.WriteMessage(websocket.TextMessage, []byte(`{"__success":"listening"}`))
			}
		} else {
			slog.Info("send", "channel", msg.Parsed.Name, "listeners", listeners)

			if !slices.Contains(listeners, msg.Parsed.Name) {
				continue
			}

			slog.Info("send", "channel", msg.Parsed.Name)

			serialized, err := json.Marshal(msg.Parsed)
			if err != nil {
				slog.Error("serialize", "error", err)
				continue
			}

			c.WriteMessage(websocket.TextMessage, serialized)
		}
	}
}
