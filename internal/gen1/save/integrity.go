package save

import (
	"crypto/sha256"
	"fmt"
)

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

// IntegrityReport contains the results of integrity checks
type IntegrityReport struct {
	IsValid       bool
	Errors        []string
	Warnings      []string
	GameVersion   GameVersion
	ChecksumValid bool
	BagValid      bool
	MoneyValid    bool
}

// CheckIntegrity performs comprehensive integrity checks on the save file
func (s *Save) CheckIntegrity() IntegrityReport {
	report := IntegrityReport{
		IsValid:     true,
		Errors:      make([]string, 0),
		Warnings:    make([]string, 0),
		GameVersion: s.DetectGameVersion(),
	}

	// 1. Checksum validation
	report.ChecksumValid = s.ValidateChecksum()
	if !report.ChecksumValid {
		report.Errors = append(report.Errors, "Invalid checksum - save may be corrupted")
		report.IsValid = false
	}

	// 2. Bag count validation
	bagCount := s.GetByte(OffsetBagCount)
	if bagCount > MaxBagItems {
		report.Errors = append(report.Errors, fmt.Sprintf("Bag count %d exceeds maximum %d", bagCount, MaxBagItems))
		report.IsValid = false
		report.BagValid = false
	} else {
		report.BagValid = true
	}

	// 3. Bag terminator check (should be 0xFF after last item)
	if bagCount < MaxBagItems {
		terminatorOffset := OffsetBagItems + (int(bagCount) * 2)
		terminator := s.GetByte(terminatorOffset)
		if terminator != 0xFF {
			report.Warnings = append(report.Warnings, "Missing bag terminator byte (0xFF)")
		}
	}

	// 4. Money validation (checked via money package)
	// Money is stored in BCD, so we check if the bytes are valid BCD
	report.MoneyValid = true // Assume valid until proven otherwise
	moneyBytes := s.GetBytes(OffsetMoney, 3)
	if moneyBytes != nil {
		for i, b := range moneyBytes {
			high := (b >> 4) & 0x0F
			low := b & 0x0F
			if high > 9 || low > 9 {
				report.Errors = append(report.Errors, fmt.Sprintf("Invalid BCD in money byte %d: 0x%02X", i, b))
				report.IsValid = false
				report.MoneyValid = false
			}
		}
	}

	// 5. Version-specific checks
	if report.GameVersion == VersionUnknown {
		report.Warnings = append(report.Warnings, "Could not detect game version - offsets may be incorrect")
	}

	return report
}

// DetectGameVersion attempts to identify the game version
func (s *Save) DetectGameVersion() GameVersion {
	// For now, we assume Yellow NA based on checksum location
	// Future: implement more sophisticated detection

	// Simple heuristic: check if checksum is at the expected offset
	if s.ValidateChecksum() {
		// Checksum is valid at 0x3523, likely Yellow NA
		return VersionYellowNA
	}

	// Check if it looks like a Yellow save by validating structure
	bagCount := s.GetByte(OffsetBagCount)
	if bagCount <= MaxBagItems {
		// Bag count is reasonable, likely Yellow NA
		return VersionYellowNA
	}

	return VersionUnknown
}

// GetSHA256 returns the SHA256 hash of the save data
func (s *Save) GetSHA256() string {
	hash := sha256.Sum256(s.data)
	return fmt.Sprintf("%x", hash)
}

// ValidateAgainstHash checks if the save matches the expected hash
func (s *Save) ValidateAgainstHash(expectedHash string) bool {
	actualHash := s.GetSHA256()
	return actualHash == expectedHash
}
