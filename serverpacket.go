package archipelago

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"maps"
	"math"
	"time"
)

// Possible values for "cmd" on packets received from the Archipelago server.
const (
	CommandRoomInfo    = "RoomInfo"
	CommandDataPackage = "DataPackage"
)

// Actions that a player can take regarding their items.
const (
	// ActionCollect collects all of a player's items from other players.
	ActionCollect = "collect"

	// ActionRelease sends all items in a player's game to other players.
	ActionRelease = "release"

	// ActionRemaining queries the number of items in a player's game.
	ActionRemaining = "remaining"
)

// Permissions for [ActionCollect], [ActionRelease], and [ActionRemaining].
const (
	// PermDisabled disables an action from running at any time.
	PermDisabled = 0

	// PermEnabled allows an action to be run manually at any time.
	PermEnabled = 1

	// PermGoal allows manual use of an action only if the player's goal has been completed.
	PermGoal = 2

	// PermAuto automatically takes the action after a player's goal has been completed.
	// Only valid for [ActionCollect] and [ActionRelease].
	PermAuto = 6

	// PermAutoEnabled allows running an action at any time,
	// and automatically runs it when a player's goal has been completed.
	// Only valid for [ActionCollect] and [ActionRelease].
	PermAutoEnabled = 7
)

type ServerPacket struct {
	Type string `json:"cmd"`
	RoomInfo
}

type RoomInfo struct {
	Version             NetworkVersion     `json:"version"`
	GeneratorVersion    NetworkVersion     `json:"generator_version"`
	Tags                []string           `json:"tags"`
	AuthRequired        bool               `json:"password"`
	Permissions         map[string]int     `json:"permissions"`
	HintCost            int                `json:"hint_cost"`
	LocationCheckPoints int                `json:"location_check_points"`
	Games               []string           `json:"games"`
	PackageChecksums    map[string]string  `json:"datapackage_checksums"`
	SeedName            string             `json:"seed_name"`
	Received            FloatUnixTimestamp `json:"time"`
}

type DataPackage map[string]GameData

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

type GameData struct {
	ItemNameToID     map[string]int `json:"item_name_to_id"`
	LocationNameToID map[string]int `json:"location_name_to_id"`
	Checksum         string         `json:"checksum"`
}

type NetworkVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Build int `json:"build"`
}

// FloatUnixTimestamp is a wrapper around [time.Time] that allows unmarshaling from a float value.
type FloatUnixTimestamp struct {
	time.Time
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
