package listener

import (
	"encoding/json"
	"log/slog"
	"strings"

	"farthergate.com/sockem/config"
	"farthergate.com/sockem/data"
	"farthergate.com/sockem/state"
	"github.com/gorilla/websocket"
)

func RecvLoop(c *websocket.Conn) {
	defer c.Close()

	authenticated := false

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			slog.Error("recv", "error", err)
			for e := state.Channels.Front(); e != nil; e = e.Next() {
				channel := e.Value.(chan data.PassedMessage)
				channel <- data.PassedMessage{
					InternalMessage: data.CLOSE,
					Conn:            c,
				}
			}
			break
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
					c.WriteMessage(mt, []byte(`{"__success":"authenticated"}`))
				} else {
					c.WriteMessage(mt, []byte(`{"error":"invalid secret key"}`))
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

		slog.Info("broadcast", "channel", parsed.Name)

		for e := state.Channels.Front(); e != nil; e = e.Next() {
			channel := e.Value.(chan data.PassedMessage)
			channel <- data.PassedMessage{
				Parsed:          parsed,
				InternalMessage: internalMsgE,
				Authenticated:   authenticated,
				Conn:            c,
			}
		}
	}
}
