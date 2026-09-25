package archipelago

import (
	"time"
)

// Possible values for "cmd" on packets received from the Archipelago server,
// indicating the type of packet being sent by the server.
const (
	TypeRoomInfo    = "RoomInfo"
	TypeDataPackage = "DataPackage"
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

type packet struct {
	Cmd string `json:"cmd"`
}

// ServerPacket is implemented by all types that can be a top-level object
// sent by the Archipelago server.
// The string returned will correspond to the "cmd" key on the server packets.
type ServerPacket interface {
	Type() string
}

// RoomInfo is the information sent in the [TypeRoomInfo] packet.
type RoomInfo struct {
	Version             Version            `json:"version"`
	GeneratorVersion    Version            `json:"generator_version"`
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

// GameData is the information contained in a [TypeDataPackage] packet.
// It contains the [DataMapping] for one or more games.
type GameData map[string]DataMapping

// DataMapping is a mapping of a game's items and locations to integer IDs.
type DataMapping struct {
	ItemNameToID     map[string]int `json:"item_name_to_id"`
	LocationNameToID map[string]int `json:"location_name_to_id"`
	Checksum         string         `json:"checksum"`
}

// Version contains version information for the randomizer generator and
// Archipelago server.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Build int `json:"build"`
}

// FloatUnixTimestamp is a wrapper around [time.Time] that allows unmarshaling from a float value.
type FloatUnixTimestamp struct {
	time.Time
}

// Type implements the [ServerPacket] interface for [RoomInfo].
func (r RoomInfo) Type() string {
	return TypeRoomInfo
}

// Type implements the [ServerPacket] interface for [GameData].
func (gd GameData) Type() string {
	return TypeDataPackage
}
