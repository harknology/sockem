package listener

import (
	"encoding/json"
	"log/slog"
	"strings"

	"farthergate.com/sockem/config"
	"farthergate.com/sockem/data"
	"github.com/gorilla/websocket"
)

func RecvLoop(c *websocket.Conn, messages chan data.PassedMessage) {
	defer c.Close()

	authenticated := false

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			slog.Error("recv", "error", err)
			continue
		}

		var parsed data.Message
		err = json.Unmarshal(message, &parsed)
		if err != nil {
			slog.Error("recv:parse", "error", err)
			continue
		}

		internalMsg := strings.TrimPrefix(parsed.Name, "__sockem:")
		internalMsgE := data.NONE
		if internalMsg != parsed.Name {
			switch internalMsg {
			case "authenticate":
				if parsed.Data == config.SECRET_KEY {
					authenticated = true
				}
				internalMsgE = data.AUTHENTICATE
				goto pushMessage // pass this on to the send loop
			case "listen":
				internalMsgE = data.LISTEN
				goto pushMessage // pass this on to the send loop
			default:
				c.WriteMessage(mt, []byte(`{"error":"unknown internal"}`))
				slog.Warn("unknown internal called", "name", internalMsg)
			}
			continue
		}

	pushMessage:
		if !authenticated {
			c.WriteMessage(mt, []byte(`{"error":"not authenticated"}`))
			continue
		}

		messages <- data.PassedMessage{
			Parsed:          parsed,
			InternalMessage: internalMsgE,
			Authenticated:   authenticated,
		}
	}
}
