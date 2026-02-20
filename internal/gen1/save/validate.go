package save

import (
	"fmt"
)

// Validate performs validation checks on the save file
func (s *Save) Validate() error {
	// Check file size
	if len(s.data) != SaveSize {
		return fmt.Errorf("invalid save file size: expected %d bytes, got %d bytes", SaveSize, len(s.data))
	}

	// Validate checksum
	if !s.ValidateChecksum() {
		return fmt.Errorf("checksum validation failed: save file may be corrupted")
	}

	return nil
}

// ValidateChecksum checks if the current checksum is correct
func (s *Save) ValidateChecksum() bool {
	calculated := s.CalculateChecksum()
	stored := s.GetByte(OffsetChecksum)
	return calculated == stored
}
