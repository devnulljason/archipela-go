package archipelago

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Run("FloatUnixTimestamp", func(t *testing.T) {
		raw := []byte(`1136239445.1234567`)
		want := FloatUnixTimestamp{time.Unix(1136239445, 123456700)} // 2006-01-02 15:04:05.1234567 -0700
		got := FloatUnixTimestamp{}

		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("error unmarshaling json: %s", err)
		}

		if !want.Equal(got.Time) {
			t.Errorf("wanted %+v but got %+v", want, got)
		}
	})
}
