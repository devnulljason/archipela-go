package archipelago

import (
	"encoding/json/jsontext"
	"fmt"
	"math"
	"time"
)

type MessageType string

const ServerRoomInfo MessageType = "RoomInfo"

type ServerPacket struct {
	Type MessageType `json:"cmd"`
	RoomInfo
}

type RoomInfo struct {
	Version             NetworkVersion               `json:"version"`
	GeneratorVersion    NetworkVersion               `json:"generator_version"`
	Tags                []string                     `json:"tags"`
	AuthRequired        bool                         `json:"password"`
	Permissions         map[PermissionKey]Permission `json:"permissions"`
	HintCost            int                          `json:"hint_cost"`
	LocationCheckPoints int                          `json:"location_check_points"`
	Games               []string                     `json:"games"`
	PackageChecksums    map[string]string            `json:"datapackage_checksums"`
	SeedName            string                       `json:"seed_name"`
	Received            FloatUnixTimestamp           `json:"time"`
}

type DataPackage struct {
	Data map[string]GameData
}

type GameData struct {
	Name             string
	ItemNameToID     map[string]int
	LocationNameToID map[string]int
	Checksum         string
}

type NetworkVersion struct {
	Major int
	Minor int
	Build int
}

func (v NetworkVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}

type Permission int

const (
	Disabled    Permission = 0
	Enabled     Permission = 1
	Goal        Permission = 2
	Auto        Permission = 6
	AutoEnabled Permission = 7
)

type PermissionKey string

const (
	Release   PermissionKey = "release"
	Collect   PermissionKey = "collect"
	Remaining PermissionKey = "remaining"
)

type FloatUnixTimestamp struct {
	time.Time
}

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
