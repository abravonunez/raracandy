package save

import (
	"testing"
)

func TestBCDConversion(t *testing.T) {
	tests := []struct {
		decimal byte
		bcd     byte
	}{
		{0, 0x00},
		{23, 0x23},
		{99, 0x99},
		{50, 0x50},
		{12, 0x12},
	}

	for _, tt := range tests {
		// Test decimal to BCD in money package
		// This is a placeholder - we'll test via money package
		t.Logf("Test case: decimal=%d, expected BCD=0x%02X", tt.decimal, tt.bcd)
	}
}

func TestChecksumCalculation(t *testing.T) {
	// Create a minimal save file (32KB)
	data := make([]byte, SaveSize)

	// Set some test data
	for i := ChecksumStart; i <= ChecksumEnd; i++ {
		data[i] = 0x01 // All bytes = 1
	}

	s := &Save{data: data}

	// Calculate expected checksum
	// Sum = (ChecksumEnd - ChecksumStart + 1) * 1
	// For our test: sum = number of bytes
	sum := byte((ChecksumEnd - ChecksumStart + 1) % 256)
	expected := ^sum // Bitwise NOT

	calculated := s.CalculateChecksum()

	if calculated != expected {
		t.Errorf("CalculateChecksum() = 0x%02X, want 0x%02X", calculated, expected)
	}
}

func TestRecalculateChecksum(t *testing.T) {
	// Create a minimal save file
	data := make([]byte, SaveSize)
	s := &Save{data: data}

	// Recalculate checksum
	s.RecalculateChecksum()

	// Verify it's stored correctly
	stored := s.GetChecksum()
	calculated := s.CalculateChecksum()

	if stored != calculated {
		t.Errorf("Stored checksum 0x%02X does not match calculated 0x%02X", stored, calculated)
	}
}
