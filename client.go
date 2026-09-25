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
	conn *websocket.Conn
	addr url.URL
}

type clientConfig struct {
	hostname string
}

// ClientOptions configure a [Client] when using the [NewClient] function.
type ClientOptions func(*clientConfig)

// WithHostname changes the hostname for the Archipelago server.
// The default hostname is "archipelago.gg".
func WithHostname(hostname string) ClientOptions {
	return func(co *clientConfig) {
		co.hostname = hostname
	}
}

// NewClient uses the provided configurations to create a websocket client to
// communicate with an Archipelago server.
// A port number for the MultiWorld room must be provided.
func NewClient(port int, opts ...ClientOptions) *Client {
	config := clientConfig{
		hostname: "archipelago.gg",
	}

	for _, optfunc := range opts {
		optfunc(&config)
	}

	client := &Client{}

	client.addr = url.URL{Scheme: "wss", Host: fmt.Sprintf("%s:%d", config.hostname, port)}
	return client
}

// Connect initiates a connection to the Archipelago server.
func (c *Client) Connect() error {
	c.addr.Scheme = "wss"
	conn, res, err := websocket.DefaultDialer.Dial(c.addr.String(), nil)
	log.Printf("%+v", res)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Close attempts to gracefully close the websocket connection to an Archipelago server.
func (c *Client) Close() error {
	return c.conn.Close()
}

// ReadPackets reads the next message sent by the Archipelago server
// and parses the packets it contains.
// The packets in the message are returned to the caller without any further processing.
// Generally there will only be a single packet per message,
// but callers should always check the length of the returned slice.
func (c *Client) ReadPackets() ([]ServerPacket, error) {
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
			return nil, fmt.Errorf("received unsupported server command: %s", comm.Cmd)
		}
	}

	return packets, nil
}

// Send sends a client package to the Archipelago server.
func (c *Client) Send(p ...GetDataPackage) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, raw)
}
