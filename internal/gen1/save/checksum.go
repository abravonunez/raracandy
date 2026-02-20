package save

// CalculateChecksum computes the checksum for the save file
// Algorithm: Sum bytes from ChecksumStart to ChecksumEnd, then apply bitwise NOT
func (s *Save) CalculateChecksum() byte {
	var sum byte = 0

	for i := ChecksumStart; i <= ChecksumEnd; i++ {
		sum += s.GetByte(i)
	}

	// Apply bitwise NOT to get the checksum
	return ^sum
}

// RecalculateChecksum updates the checksum in the save data
func (s *Save) RecalculateChecksum() {
	checksum := s.CalculateChecksum()
	s.SetByte(OffsetChecksum, checksum)
}

// GetChecksum returns the currently stored checksum
func (s *Save) GetChecksum() byte {
	return s.GetByte(OffsetChecksum)
}
