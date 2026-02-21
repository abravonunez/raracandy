package profile

// GameVersion represents the detected game version
type GameVersion int

const (
	VersionUnknown GameVersion = iota
	VersionYellowNA
	VersionYellowJP
	VersionYellowEU
	VersionRedBlueNA
)

func (v GameVersion) String() string {
	switch v {
	case VersionYellowNA:
		return "Pokémon Yellow (North America)"
	case VersionYellowJP:
		return "Pokémon Yellow (Japan)"
	case VersionYellowEU:
		return "Pokémon Yellow (Europe)"
	case VersionRedBlueNA:
		return "Pokémon Red/Blue (North America)"
	default:
		return "Unknown"
	}
}

// GameProfile encapsulates all version-specific configurations and offsets
type GameProfile struct {
	Version        GameVersion
	Name           string
	OffsetChecksum int
	ChecksumStart  int
	ChecksumEnd    int
	OffsetBagCount int
	OffsetBagItems int
	OffsetMoney    int
	MaxBagItems    int
	MaxMoney       uint32
}

var (
	// ProfileYellowNA defines offsets and config for Pokémon Yellow (North America)
	ProfileYellowNA = &GameProfile{
		Version:        VersionYellowNA,
		Name:           "Pokémon Yellow (North America)",
		OffsetChecksum: 0x3523,
		ChecksumStart:  0x2598,
		ChecksumEnd:    0x3522,
		OffsetBagCount: 0x25C9,
		OffsetBagItems: 0x25CA,
		OffsetMoney:    0x25F3,
		MaxBagItems:    20,
		MaxMoney:       999999,
	}

	// ProfileRedBlueNA defines offsets and config for Pokémon Red/Blue (North America)
	ProfileRedBlueNA = &GameProfile{
		Version:        VersionRedBlueNA,
		Name:           "Pokémon Red/Blue (North America)",
		OffsetChecksum: 0x3523,
		ChecksumStart:  0x2598,
		ChecksumEnd:    0x3522,
		OffsetBagCount: 0x25C9,
		OffsetBagItems: 0x25CA,
		OffsetMoney:    0x25F3,
		MaxBagItems:    20,
		MaxMoney:       999999,
	}
)

// GetProfile returns the appropriate GameProfile for a given game version
func GetProfile(version GameVersion) *GameProfile {
	switch version {
	case VersionYellowNA:
		return ProfileYellowNA
	case VersionRedBlueNA:
		return ProfileRedBlueNA
	default:
		// Fallback to Yellow NA as a conservative default
		return ProfileYellowNA
	}
}
