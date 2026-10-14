package archipelago

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		port int
		opts []ClientOptions
		want *Client
	}{
		{
			name: "default hostname",
			port: 12345,
			want: &Client{
				addr: url.URL{
					Scheme: "wss",
					Host:   "archipelago.gg:12345",
				},
			},
		},
		{
			name: "custom hostname",
			port: 12345,
			opts: []ClientOptions{WithHostname("foo")},
			want: &Client{
				addr: url.URL{
					Scheme: "wss",
					Host:   "foo:12345",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.port, tt.opts...)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %+v but got %+v", tt.want, got)
			}
		})
	}
}

func TestReceive(t *testing.T) {
	tests := []struct {
		name      string
		rawPacket string
		want      ServerPacket
	}{
		{
			name: "RoomInfo",
			rawPacket: `
[{
	"cmd": "RoomInfo",
	"password": false,
	"games": ["Archipelago"],
	"tags": ["TestClient"],
	"version": { "major": 0, "minor": 6, "build": 8, "class": "Version" },
	"generator_version": {"major": 0, "minor": 6, "build":7, "class": "Version" },
	"permissions": { "release": 7, "remaining": 2,"collect": 1 },
	"hint_cost": 20,
	"location_check_points": 1,
	"datapackage_checksums": { "Archipelago":"ac9141e9ad0318df2fa27da5f20c50a842afeecb" },
	"seed_name": "83353347295836478478",
	"time": 1791925809.8723803
}]`,
			want: RoomInfo{
				Version:             version{Major: 0, Minor: 6, Build: 8},
				GeneratorVersion:    version{Major: 0, Minor: 6, Build: 7},
				Tags:                []string{"TestClient"},
				AuthRequired:        false,
				Permissions:         map[string]int{"release": 7, "remaining": 2, "collect": 1},
				HintCost:            20,
				LocationCheckPoints: 1,
				Games:               []string{"Archipelago"},
				PackageChecksums:    map[string]string{"Archipelago": "ac9141e9ad0318df2fa27da5f20c50a842afeecb"},
				SeedName:            "83353347295836478478",
				Received:            floatUnixTimestamp{time.Unix(1791925809, 872380300)},
			},
		},
		{
			name: "GameData",
			rawPacket: `
[{
	"cmd": "DataPackage",
	"data": {
		"games": {
			"Archipelago": {
				"item_name_to_id": {"Nothing": -1},
				"location_name_to_id": {"Cheat Console": -1, "Server": -2},
				"checksum": "ac9141e9ad0318df2fa27da5f20c50a842afeecb"
			}
		}
	}
}]`,
			want: GameData{
				"Archipelago": {
					ItemID: map[string]int{"Nothing": -1},
					LocationID: map[string]int{
						"Cheat Console": -1,
						"Server":        -2,
					},
					Checksum: "ac9141e9ad0318df2fa27da5f20c50a842afeecb",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				conn: &mockConnection{
					raw: tt.rawPacket,
				},
			}
			msgs, err := client.Receive()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(msgs) != 1 {
				t.Fatalf("expected 1 message but got %+v", msgs)
			}

			got := msgs[0]
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Fatalf("mismatched types: wanted %T but got %T", tt.want, got)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("wanted %+v but got %+v", tt.want, got)
			}
		})
	}
}

type mockConnection struct {
	raw string
}

func (m *mockConnection) Close() error {
	return nil
}

func (m mockConnection) ReadMessage() (int, []byte, error) {
	return 0, []byte(m.raw), nil
}

func (m *mockConnection) WriteMessage(messageType int, data []byte) error {
	return nil
}
