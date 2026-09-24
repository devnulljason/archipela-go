package archipelago

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"maps"
	"math"
	"time"
)

// UnmarshalJSONFrom implements the [encoding/json/v2.UnmarshalerFrom] interface.
func (dp *DataPackage) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	type dataPackage struct {
		Data struct {
			Games map[string]GameData `json:"games"`
		} `json:"data"`
	}

	d := dataPackage{}
	if err := json.UnmarshalDecode(dec, &d); err != nil {
		return err
	}

	maps.Copy(*dp, d.Data.Games)
	return nil
}

// UnmarshalJSONFrom implements the [encoding/json/v2.UnmarshalerFrom] interface.
func (u *FloatUnixTimestamp) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	timeToken, err := dec.ReadToken()
	if err != nil {
		return err
	}
	timeFloat, err := timeToken.Float()
	if err != nil {
		return err
	}
	sec, nsec := math.Modf(timeFloat)

	// Python floating-point timestamps don't have full nanosecond precision, see PEP-564
	// converting in two steps to avoid floating-point shenanigans
	intnsec := int64(nsec*1e7) * 100
	u.Time = time.Unix(int64(sec), intnsec)
	return nil
}
