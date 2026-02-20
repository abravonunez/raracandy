package save

import (
	"fmt"
	"os"
)

const (
	// Save file constants for Pokemon Yellow (NA)
	SaveSize = 0x8000 // 32 KB
	BankSize = 0x2000 // 8 KB per bank

	// Checksum offsets
	OffsetChecksum = 0x3523
	ChecksumStart  = 0x2598
	ChecksumEnd    = 0x3522

	// Item bag offsets (estimated, need verification)
	OffsetBagCount = 0x25C9 // D31D - AD54
	OffsetBagItems = 0x25CA // D31E - AD54
	MaxBagItems    = 20

	// Money offsets (estimated, need verification)
	OffsetMoney = 0x25F3 // D347 - AD54
	MaxMoney    = 999999

	// Item IDs
	ItemRareCandy = 0x28
)

// Save represents a Pokemon Gen 1 save file
type Save struct {
	data     []byte
	filePath string
}

// Load reads a save file from disk
func Load(path string) (*Save, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	s := &Save{
		data:     data,
		filePath: path,
	}

	// Validate the save file
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("save file validation failed: %w", err)
	}

	return s, nil
}

// Write saves the data to a file
func (s *Save) Write(path string) error {
	// Recalculate checksum before writing
	s.RecalculateChecksum()

	// Write to file with proper permissions
	if err := os.WriteFile(path, s.data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// Data returns a copy of the save data
func (s *Save) Data() []byte {
	cpy := make([]byte, len(s.data))
	copy(cpy, s.data)
	return cpy
}

// GetByte returns the byte at the given offset
func (s *Save) GetByte(offset int) byte {
	if offset < 0 || offset >= len(s.data) {
		return 0
	}
	return s.data[offset]
}

// SetByte sets the byte at the given offset
func (s *Save) SetByte(offset int, value byte) error {
	if offset < 0 || offset >= len(s.data) {
		return fmt.Errorf("offset %d out of bounds (size: %d)", offset, len(s.data))
	}
	s.data[offset] = value
	return nil
}

// GetBytes returns a slice of bytes starting at offset with given length
func (s *Save) GetBytes(offset, length int) []byte {
	if offset < 0 || offset+length > len(s.data) {
		return nil
	}
	result := make([]byte, length)
	copy(result, s.data[offset:offset+length])
	return result
}

// SetBytes sets multiple bytes starting at offset
func (s *Save) SetBytes(offset int, data []byte) error {
	if offset < 0 || offset+len(data) > len(s.data) {
		return fmt.Errorf("offset %d + length %d out of bounds (size: %d)", offset, len(data), len(s.data))
	}
	copy(s.data[offset:], data)
	return nil
}
