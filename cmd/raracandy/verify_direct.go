package main

import "github.com/spf13/cobra"

var verifyDirectCmd = &cobra.Command{
	Use:   "verify <save-file>",
	Short: "Verify save file integrity and structure (auto-detects version)",
	Long: `Performs comprehensive integrity checks on a save file without modifying it.
Automatically detects whether the save file is from Pokémon Red, Blue, or Yellow.

Checks include:
- File size validation
- Checksum verification
- Game version detection
- Bag structure validation
- Money format validation
- SHA256 hash (optional)`,
	Args: cobra.ExactArgs(1),
	RunE: runVerify, // Reuse the same logic from verify.go
}

func init() {
	rootCmd.AddCommand(verifyDirectCmd)

	verifyDirectCmd.Flags().StringVar(&verifyExpectedHash, "expected-hash", "", "Expected SHA256 hash to verify against")
}
