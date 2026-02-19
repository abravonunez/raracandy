package items

import (
	"fmt"
	"strings"
)

// Item ID constants for Pokemon Yellow
const (
	IDRareCandy     = 0x28
	IDMasterBall    = 0x01
	IDUltraBall     = 0x02
	IDGreatBall     = 0x03
	IDPokeBall      = 0x04
	IDPotion        = 0x14
	IDSuperPotion   = 0x15
	IDHyperPotion   = 0x16
	IDMaxPotion     = 0x17
	IDFullRestore   = 0x18
	IDRevive        = 0x19
	IDMaxRevive     = 0x1A
	IDEscape        = 0x1C
	IDRepel         = 0x1D
	IDSuperRepel    = 0x1E
	IDMaxRepel      = 0x1F
	IDAntidote      = 0x20
	IDBurnHeal      = 0x21
	IDIceHeal       = 0x22
	IDAwakening     = 0x23
	IDParalyzeHeal  = 0x24
	IDFullHeal      = 0x25
)

// itemNames maps item IDs to human-readable names
var itemNames = map[byte]string{
	IDRareCandy:     "Rare Candy",
	IDMasterBall:    "Master Ball",
	IDUltraBall:     "Ultra Ball",
	IDGreatBall:     "Great Ball",
	IDPokeBall:      "Poké Ball",
	IDPotion:        "Potion",
	IDSuperPotion:   "Super Potion",
	IDHyperPotion:   "Hyper Potion",
	IDMaxPotion:     "Max Potion",
	IDFullRestore:   "Full Restore",
	IDRevive:        "Revive",
	IDMaxRevive:     "Max Revive",
	IDEscape:        "Escape Rope",
	IDRepel:         "Repel",
	IDSuperRepel:    "Super Repel",
	IDMaxRepel:      "Max Repel",
	IDAntidote:      "Antidote",
	IDBurnHeal:      "Burn Heal",
	IDIceHeal:       "Ice Heal",
	IDAwakening:     "Awakening",
	IDParalyzeHeal:  "Paralyze Heal",
	IDFullHeal:      "Full Heal",
}

// itemIDs maps human-readable names to item IDs (lowercase for case-insensitive lookup)
var itemIDs = map[string]byte{
	"rare_candy":    IDRareCandy,
	"rarecandy":     IDRareCandy,
	"master_ball":   IDMasterBall,
	"masterball":    IDMasterBall,
	"ultra_ball":    IDUltraBall,
	"ultraball":     IDUltraBall,
	"great_ball":    IDGreatBall,
	"greatball":     IDGreatBall,
	"poke_ball":     IDPokeBall,
	"pokeball":      IDPokeBall,
	"potion":        IDPotion,
	"super_potion":  IDSuperPotion,
	"superpotion":   IDSuperPotion,
	"hyper_potion":  IDHyperPotion,
	"hyperpotion":   IDHyperPotion,
	"max_potion":    IDMaxPotion,
	"maxpotion":     IDMaxPotion,
	"full_restore":  IDFullRestore,
	"fullrestore":   IDFullRestore,
	"revive":        IDRevive,
	"max_revive":    IDMaxRevive,
	"maxrevive":     IDMaxRevive,
	"escape_rope":   IDEscape,
	"escaperope":    IDEscape,
	"repel":         IDRepel,
	"super_repel":   IDSuperRepel,
	"superrepel":    IDSuperRepel,
	"max_repel":     IDMaxRepel,
	"maxrepel":      IDMaxRepel,
	"antidote":      IDAntidote,
	"burn_heal":     IDBurnHeal,
	"burnheal":      IDBurnHeal,
	"ice_heal":      IDIceHeal,
	"iceheal":       IDIceHeal,
	"awakening":     IDAwakening,
	"paralyze_heal": IDParalyzeHeal,
	"paralyzeheal":  IDParalyzeHeal,
	"full_heal":     IDFullHeal,
	"fullheal":      IDFullHeal,
}

// GetItemID returns the item ID for a given name (case-insensitive)
func GetItemID(name string) (byte, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	id, ok := itemIDs[normalized]
	if !ok {
		return 0, fmt.Errorf("unknown item: %s", name)
	}
	return id, nil
}

// GetItemName returns the human-readable name for an item ID
func GetItemName(id byte) string {
	name, ok := itemNames[id]
	if !ok {
		return fmt.Sprintf("Unknown Item (0x%02X)", id)
	}
	return name
}

// IsValidItemID checks if an item ID is recognized
func IsValidItemID(id byte) bool {
	_, ok := itemNames[id]
	return ok
}
