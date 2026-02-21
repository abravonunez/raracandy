package save

import (
	"crypto/sha256"
	"fmt"

	"github.com/abravonunez/raracandy/internal/gen1/profile"
)

// IntegrityReport contains the results of integrity checks
type IntegrityReport struct {
	IsValid       bool
	Errors        []string
	Warnings      []string
	GameVersion   profile.GameVersion
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
	prof := s.GetProfile()
	bagCount := s.GetByte(prof.OffsetBagCount)
	if bagCount > byte(prof.MaxBagItems) {
		report.Errors = append(report.Errors, fmt.Sprintf("Bag count %d exceeds maximum %d", bagCount, prof.MaxBagItems))
		report.IsValid = false
		report.BagValid = false
	} else {
		report.BagValid = true
	}

	// 3. Bag terminator check (should be 0xFF after last item)
	if bagCount < byte(prof.MaxBagItems) {
		terminatorOffset := prof.OffsetBagItems + (int(bagCount) * 2)
		terminator := s.GetByte(terminatorOffset)
		if terminator != 0xFF {
			report.Warnings = append(report.Warnings, "Missing bag terminator byte (0xFF)")
		}
	}

	// 4. Money validation (checked via money package)
	// Money is stored in BCD, so we check if the bytes are valid BCD
	report.MoneyValid = true // Assume valid until proven otherwise
	moneyBytes := s.GetBytes(prof.OffsetMoney, 3)
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
	if report.GameVersion == profile.VersionUnknown {
		report.Warnings = append(report.Warnings, "Could not detect game version - offsets may be incorrect")
	}

	return report
}

// DetectGameVersion attempts to identify the game version
// This uses hardcoded NA offsets for detection since all NA versions share the same offsets
func (s *Save) DetectGameVersion() profile.GameVersion {
	// Constants for NA version detection (same for Red/Blue/Yellow NA)
	const (
		naOffsetChecksum = 0x3523
		naChecksumStart  = 0x2598
		naChecksumEnd    = 0x3522
		naOffsetBagCount = 0x25C9
		naOffsetMoney    = 0x25F3
	)

	// 1. Validate checksum at NA offset
	var sum byte = 0
	for i := naChecksumStart; i <= naChecksumEnd; i++ {
		sum += s.GetByte(i)
	}
	calculatedChecksum := ^sum
	storedChecksum := s.GetByte(naOffsetChecksum)

	checksumValid := (calculatedChecksum == storedChecksum)

	// 2. Validate bag structure
	bagCount := s.GetByte(naOffsetBagCount)
	bagValid := bagCount <= 20

	// 3. Validate money (BCD format)
	moneyValid := true
	moneyBytes := s.GetBytes(naOffsetMoney, 3)
	if moneyBytes != nil {
		for _, b := range moneyBytes {
			high := (b >> 4) & 0x0F
			low := b & 0x0F
			if high > 9 || low > 9 {
				moneyValid = false
				break
			}
		}
	}

	// If checksum is valid and structure looks good, it's a NA version
	if checksumValid && bagValid && moneyValid {
		// Default to Yellow NA since offsets are identical for Red/Blue/Yellow NA
		// Future: could differentiate by checking for Pikachu-specific data
		return profile.VersionYellowNA
	}

	// If structure looks reasonable but checksum is invalid, still assume NA
	// (checksum might be invalid for legitimate reasons, e.g., corrupted save)
	if bagValid && moneyValid {
		return profile.VersionYellowNA
	}

	return profile.VersionUnknown
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
