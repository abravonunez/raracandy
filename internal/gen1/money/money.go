package money

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/gen1/save"
)

const (
	MaxMoney = 999999
)

// GetMoney reads and decodes the player's money from the save file
// Money is stored as 3 bytes in BCD (Binary-Coded Decimal) format
func GetMoney(s *save.Save) uint32 {
	profile := s.GetProfile()
	bytes := s.GetBytes(profile.OffsetMoney, 3)
	if bytes == nil || len(bytes) != 3 {
		return 0
	}

	// Decode BCD: each byte represents two decimal digits
	// Example: 0x12 0x34 0x56 = 123,456
	money := uint32(bcdToDecimal(bytes[0]))*10000 +
		uint32(bcdToDecimal(bytes[1]))*100 +
		uint32(bcdToDecimal(bytes[2]))

	return money
}

// SetMoney encodes and writes the player's money to the save file
// Returns error if amount exceeds MaxMoney
func SetMoney(s *save.Save, amount uint32) error {
	if amount > MaxMoney {
		return fmt.Errorf("amount %d exceeds maximum %d", amount, MaxMoney)
	}

	// Encode to BCD
	bytes := make([]byte, 3)
	bytes[0] = decimalToBCD(byte(amount / 10000))     // Ten thousands and thousands
	bytes[1] = decimalToBCD(byte((amount / 100) % 100)) // Hundreds and tens
	bytes[2] = decimalToBCD(byte(amount % 100))        // Units

	profile := s.GetProfile()
	return s.SetBytes(profile.OffsetMoney, bytes)
}

// bcdToDecimal converts a BCD byte to decimal
// Example: 0x23 -> 23
func bcdToDecimal(bcd byte) byte {
	high := (bcd >> 4) & 0x0F // High nibble
	low := bcd & 0x0F          // Low nibble
	return high*10 + low
}

// decimalToBCD converts a decimal byte (0-99) to BCD
// Example: 23 -> 0x23
func decimalToBCD(dec byte) byte {
	if dec > 99 {
		dec = 99
	}
	high := dec / 10
	low := dec % 10
	return (high << 4) | low
}

// FormatMoney returns a formatted money string with thousands separator
// Example: 123456 -> "¥123,456"
func FormatMoney(amount uint32) string {
	// Format with thousands separator
	if amount >= 1000000 {
		return fmt.Sprintf("¥%d,%03d,%03d", amount/1000000, (amount/1000)%1000, amount%1000)
	} else if amount >= 1000 {
		return fmt.Sprintf("¥%d,%03d", amount/1000, amount%1000)
	}
	return fmt.Sprintf("¥%d", amount)
}
