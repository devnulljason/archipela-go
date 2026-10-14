package archipelago

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

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
