package archipelago

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

type GetDataPackage struct {
	Games []string
}

func (p *GetDataPackage) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("cmd")); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("GetDataPackage")); err != nil {
		return err
	}

	if len(p.Games) > 0 {
		if err := enc.WriteToken(jsontext.String("games")); err != nil {
			return err
		}

		if err := json.MarshalEncode(enc, p.Games); err != nil {
			return err
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}
