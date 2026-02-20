package save

// CreateTestSave creates a minimal valid Pokemon Yellow save file for testing
func CreateTestSave() *Save {
	data := make([]byte, SaveSize)

	// Initialize with zeros
	for i := range data {
		data[i] = 0x00
	}

	s := &Save{
		data:     data,
		filePath: "test.sav",
	}

	// Set up minimal bag (empty for now)
	s.SetByte(OffsetBagCount, 0)
	s.SetByte(OffsetBagItems, 0xFF) // Terminator

	// Set money to 0
	s.SetBytes(OffsetMoney, []byte{0x00, 0x00, 0x00})

	// Calculate and set checksum
	s.RecalculateChecksum()

	return s
}
