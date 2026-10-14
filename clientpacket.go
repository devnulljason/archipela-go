package archipelago

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// GetDataPackage represents a request for metadata about the Archipelago room's games.
// The metadata is provided in a [GameData] response from the server.
type GetDataPackage struct {
	Games []string `json:"games"`
}

// MarshalJSONTo implements the [encoding/json/v2.MarshalerTo] interface for [GetDataPackage].
func (gdp GetDataPackage) MarshalJSONTo(enc *jsontext.Encoder) error {
	type dataPackage GetDataPackage
	return json.MarshalEncode(enc, struct {
		dataPackage

		Cmd string `json:"cmd"`
	}{
		Cmd:         "GetDataPackage",
		dataPackage: dataPackage(gdp),
	},
		json.FormatNilSliceAsNull(false))
}
