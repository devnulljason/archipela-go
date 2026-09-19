package archipelago

import (
	"encoding/json/v2"
	"testing"
)

func TestJSONMarshal(t *testing.T) {
	testcases := []struct {
		name string
		obj  any
		want string
	}{
		{
			name: "GetDataPackage",
			obj:  GetDataPackage{Games: []string{"Archipelago"}},
			want: `{"cmd":"GetDataPackage","games":["Archipelago"]}`,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.obj)
			if err != nil {
				t.Fatalf("error marshaling json: %s", err)
			}

			if string(got) != tt.want {
				t.Errorf("wanted %s but got %s", tt.want, string(got))
			}
		})
	}
}
