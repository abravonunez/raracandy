package save

// CalculateChecksum computes the checksum for the save file
// Algorithm: Sum bytes from ChecksumStart to ChecksumEnd, then apply bitwise NOT
func (s *Save) CalculateChecksum() byte {
	var sum byte = 0

	profile := s.GetProfile()
	for i := profile.ChecksumStart; i <= profile.ChecksumEnd; i++ {
		sum += s.GetByte(i)
	}

	// Apply bitwise NOT to get the checksum
	return ^sum
}

// RecalculateChecksum updates the checksum in the save data
func (s *Save) RecalculateChecksum() {
	checksum := s.CalculateChecksum()
	profile := s.GetProfile()
	s.SetByte(profile.OffsetChecksum, checksum)
}

// GetChecksum returns the currently stored checksum
func (s *Save) GetChecksum() byte {
	profile := s.GetProfile()
	return s.GetByte(profile.OffsetChecksum)
}
