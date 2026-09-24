// Package archipelago implements the [Archipelago network protocol].
//
// [Archipelago network protocol]: https://github.com/ArchipelagoMW/Archipelago/blob/0.6.7/docs/network%20protocol.md
package archipelago

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

type clientOptions struct {
	hostname string
}

type ClientOption func(*clientOptions)

func WithHostname(hostname string) ClientOption {
	return func(co *clientOptions) {
		co.hostname = hostname
	}
}

func NewClient(port int, opts ...ClientOption) (*Client, error) {
	defaults := &clientOptions{
		hostname: "archipelago.gg",
	}

	for _, opt := range opts {
		opt(defaults)
	}

	addr := url.URL{Scheme: "wss", Host: fmt.Sprintf("%s:%d", defaults.hostname, port)}
	conn, _, err := websocket.DefaultDialer.Dial(addr.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn: conn,
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Receive() ([]ServerPacket, error) {
	_, raw, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(raw))

	packets := []ServerPacket{}
	if err = json.Unmarshal(raw, &packets); err != nil {
		return nil, err
	}

	return packets, nil
}

func (c *Client) Send(p ...GetDataPackage) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, raw)
}
