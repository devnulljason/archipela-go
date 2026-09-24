package archipelago

import (
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

type Packet struct {
	Command string `json:"cmd"`
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

type GameData map[string]DataMapping

type DataMapping struct {
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
