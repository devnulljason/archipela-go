// Package archipelago implements the [Archipelago network protocol].
//
// [Archipelago network protocol]: https://github.com/ArchipelagoMW/Archipelago/blob/0.6.7/docs/network%20protocol.md
package archipelago

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/gorilla/websocket"
)

func Connect(host string, port int) (err error) {
	hostport := fmt.Sprintf("%s:%d", host, port)
	addr := url.URL{Scheme: "wss", Host: hostport}
	slog.Debug("connecting", "address", addr.String())

	conn, _, err := websocket.DefaultDialer.Dial(addr.String(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			err = closeErr
			slog.Error("closing connection", "error", err)
		}
	}()

	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		slog.Error("reading message", "error", err)
		return err
	}
	if msgType != websocket.TextMessage {
		err = errors.New("unexpected message type")
		slog.Error(err.Error(), "type", msgType)
		return err
	}
	slog.Info("received", "msg", msg)

	return err
}
