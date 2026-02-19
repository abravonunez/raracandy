package main

import (
	"fmt"
	"os"

	"github.com/abravonunez/raracandy/internal/gen1/save"
)

func main() {
	// Create a test save
	s := save.CreateTestSave()

	// Write to file
	outputPath := "testdata/fixtures/test.sav"
	if err := s.Write(outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Test save created: %s\n", outputPath)
	fmt.Printf("Size: %d bytes\n", len(s.Data()))
	fmt.Printf("Checksum: 0x%02X\n", s.GetChecksum())
}
