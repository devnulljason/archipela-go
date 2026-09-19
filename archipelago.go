// Package archipelago implements the [Archipelago network protocol].
//
// [Archipelago network protocol]: https://github.com/ArchipelagoMW/Archipelago/blob/0.6.7/docs/network%20protocol.md
package archipelago

import (
	"encoding/json"
	"errors"
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

func (c *Client) Receive() error {
	_, raw, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	packets := []ServerPacket{}
	if err = json.Unmarshal(raw, &packets); err != nil {
		return err
	}

	for _, p := range packets {
		switch p.Type {
		case ServerRoomInfo:
			fmt.Printf("Generator version: %s\n", p.GeneratorVersion)
			fmt.Printf("Archipelago version: %s\n", p.Version)
			return nil
		default:
			fmt.Println(errors.New("unknown packet type"))
		}
	}
	return nil
}
