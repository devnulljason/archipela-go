package archipelago

import (
	"encoding/json/v2"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Run("DataPackage", func(t *testing.T) {
		raw := `{
    "cmd": "DataPackage",
    "data": {
        "games": {
            "Archipelago": {
                "item_name_to_id": {
                    "Nothing": -1
                },
                "location_name_to_id": {
                    "Cheat Console": -1,
                    "Server": -2
                },
                "checksum": "ac9141e9ad0318df2fa27da5f20c50a842afeecb"
            }
        }
    }
}`
		want := GameData{
			"Archipelago": {
				ItemNameToID: map[string]int{
					"Nothing": -1,
				},
				LocationNameToID: map[string]int{
					"Cheat Console": -1,
					"Server":        -2,
				},
				Checksum: "ac9141e9ad0318df2fa27da5f20c50a842afeecb",
			},
		}
		got := GameData{}

		if err := json.Unmarshal([]byte(raw), &got); err != nil {
			t.Fatalf("error unmarshaling json: %s", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wanted %+v but got %+v", want, got)
		}
	})

	t.Run("FloatUnixTimestamp", func(t *testing.T) {
		raw := []byte("1136239445.1234567")
		want := FloatUnixTimestamp{time.Unix(1136239445, 123456700)}
		got := FloatUnixTimestamp{}

		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("error unmarshaling json: %s", err)
		}

		if !want.Equal(got.Time) {
			t.Errorf("wanted %+v but got %+v", want, got)
		}
	})
}
