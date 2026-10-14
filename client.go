// Package archipelago implements the [Archipelago network protocol].
//
// [Archipelago network protocol]: https://github.com/ArchipelagoMW/Archipelago/blob/0.6.7/docs/network%20protocol.md
package archipelago

import (
	"encoding/json"
	"encoding/json/jsontext"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client provides a websocket client for communicating with an Archipelago server.
type Client struct {
	conn wsConnection
	addr url.URL
}

// ClientOptions are configurations for a [Client] that can be used with the [NewClient] function.
type ClientOptions func(*clientConfig)

type clientConfig struct {
	hostname string
}

type wsConnection interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	Close() error
}

// WithHostname changes the hostname for the Archipelago server.
// The default hostname is "archipelago.gg".
func WithHostname(hostname string) ClientOptions {
	return func(co *clientConfig) {
		co.hostname = hostname
	}
}

// NewClient uses the provided configurations to create a websocket client for
// communication with an Archipelago server.
// A port number for the MultiWorld room must be provided.
// The client waits to establish a connection with a server until the first read or write.
func NewClient(port int, opts ...ClientOptions) *Client {
	config := clientConfig{
		hostname: "archipelago.gg",
	}

	for _, opt := range opts {
		opt(&config)
	}

	client := Client{
		addr: url.URL{Scheme: "wss", Host: fmt.Sprintf("%s:%d", config.hostname, port)},
	}

	return &client
}

// Close attempts to gracefully close the websocket connection to an Archipelago server.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Receive reads the next message sent by the Archipelago server
// and parses the packets it contains into Go structs.
// The packets in the message are returned to the caller without any further processing.
// Most of the time there will only be a single packet per message,
// but callers should always check the length of the returned slice.
func (c *Client) Receive() ([]ServerPacket, error) {
	if err := c.ensureConnection(); err != nil {
		return nil, err
	}

	_, raw, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	data := []jsontext.Value{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	packets := make([]ServerPacket, len(data))
	for i, d := range data {
		comm := packet{}
		if err := json.Unmarshal(d, &comm); err != nil {
			return nil, err
		}

		switch comm.Cmd {
		case TypeDataPackage:
			p := GameData{}
			if err := json.Unmarshal(d, &p); err != nil {
				return nil, err
			}
			packets[i] = p
		case TypeRoomInfo:
			p := RoomInfo{}
			if err := json.Unmarshal(d, &p); err != nil {
				return nil, err
			}
			packets[i] = p
		default:
			return nil, fmt.Errorf("unsupported server message: %s", comm.Cmd)
		}
	}

	return packets, nil
}

// Send sends a client message to the Archipelago server.
func (c Client) Send(p ...GetDataPackage) error {
	if err := c.ensureConnection(); err != nil {
		return err
	}

	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, raw)
}

func (c *Client) ensureConnection() error {
	if c.conn != nil {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.addr.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}
