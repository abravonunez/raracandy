package main

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/gen1/profile"
	"github.com/abravonunez/raracandy/internal/gen1/save"
	"github.com/spf13/cobra"
)

var (
	verifyExpectedHash string
)

var verifyCmd = &cobra.Command{
	Use:   "verify <save-file>",
	Short: "Verify save file integrity and structure",
	Long: `Performs comprehensive integrity checks on a save file without modifying it.
Checks include:
- File size validation
- Checksum verification
- Game version detection
- Bag structure validation
- Money format validation
- SHA256 hash (optional)`,
	Args: cobra.ExactArgs(1),
	RunE: runVerify,
}

func init() {
	yellowCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVar(&verifyExpectedHash, "expected-hash", "", "Expected SHA256 hash to verify against")
}

func runVerify(cmd *cobra.Command, args []string) error {
	savePath := args[0]

	// Load save file
	s, err := save.Load(savePath)
	if err != nil {
		return fmt.Errorf("failed to load save: %w", err)
	}

	fmt.Printf("Save File: %s\n", savePath)
	fmt.Printf("Size: %d KB\n", len(s.Data())/1024)
	fmt.Println()

	// Run integrity check
	report := s.CheckIntegrity()

	// Game version
	fmt.Printf("Detected Version: %s\n", report.GameVersion)
	if report.GameVersion == profile.VersionUnknown {
		fmt.Println("  ⚠️  Warning: Unknown version - offsets may be incorrect")
	}
	fmt.Println()

	// Checksum
	fmt.Println("Checksum:")
	fmt.Printf("  Stored:     0x%02X\n", s.GetChecksum())
	fmt.Printf("  Calculated: 0x%02X\n", s.CalculateChecksum())
	if report.ChecksumValid {
		fmt.Println("  Status:     ✓ Valid")
	} else {
		fmt.Println("  Status:     ✗ Invalid")
	}
	fmt.Println()

	// Bag validation
	fmt.Println("Bag Structure:")
	if report.BagValid {
		fmt.Println("  Status:     ✓ Valid")
	} else {
		fmt.Println("  Status:     ✗ Invalid")
	}
	fmt.Println()

	// Money validation
	fmt.Println("Money Format:")
	if report.MoneyValid {
		fmt.Println("  Status:     ✓ Valid BCD encoding")
	} else {
		fmt.Println("  Status:     ✗ Invalid BCD encoding")
	}
	fmt.Println()

	// SHA256 hash
	hash := s.GetSHA256()
	fmt.Printf("SHA256: %s\n", hash)

	if verifyExpectedHash != "" {
		if s.ValidateAgainstHash(verifyExpectedHash) {
			fmt.Println("  Status:     ✓ Hash matches expected value")
		} else {
			fmt.Println("  Status:     ✗ Hash does NOT match")
			return fmt.Errorf("hash mismatch")
		}
	}
	fmt.Println()

	// Errors
	if len(report.Errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range report.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
		fmt.Println()
	}

	// Warnings
	if len(report.Warnings) > 0 {
		fmt.Println("Warnings:")
		for _, warn := range report.Warnings {
			fmt.Printf("  ⚠️  %s\n", warn)
		}
		fmt.Println()
	}

	// Overall status
	if report.IsValid {
		fmt.Println("Overall Status: ✓ VALID - Safe to modify")
		return nil
	} else {
		fmt.Println("Overall Status: ✗ INVALID - Do NOT modify this save!")
		return fmt.Errorf("integrity check failed")
	}
}
