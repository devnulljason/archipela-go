package archipelago

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"maps"
	"strconv"
	"strings"
	"time"
)

// UnmarshalJSONFrom implements the [encoding/json/v2.UnmarshalerFrom] interface.
func (gd *GameData) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	type dataPackage struct {
		Data struct {
			Games map[string]dataMapping `json:"games"`
		} `json:"data"`
	}

	d := dataPackage{}
	if err := json.UnmarshalDecode(dec, &d); err != nil {
		return err
	}

	maps.Copy(*gd, d.Data.Games)
	return nil
}

// UnmarshalJSONFrom implements the [encoding/json/v2.UnmarshalerFrom] interface.
func (u *floatUnixTimestamp) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	token, err := dec.ReadToken()
	if err != nil {
		return err
	}
	// Python floating-point timestamps don't have full nanosecond precision, see PEP-564
	// using strings to avoid floating-point shenanigans
	timeStr := strings.Split(token.String(), ".")
	sec, err := strconv.Atoi(timeStr[0])
	if err != nil {
		return err
	}
	nsec, err := strconv.Atoi(timeStr[1])
	if err != nil {
		return err
	}

	u.Time = time.Unix(int64(sec), int64(nsec*100))
	return nil
}
