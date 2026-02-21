package main

import "github.com/spf13/cobra"

var inspectDirectCmd = &cobra.Command{
	Use:   "inspect <save-file>",
	Short: "Inspect a Pokémon Gen 1 save file (auto-detects version)",
	Long: `Display information about a save file including money, bag items, and checksum status.
Automatically detects whether the save file is from Pokémon Red, Blue, or Yellow.`,
	Args: cobra.ExactArgs(1),
	RunE: runInspect, // Reuse the same logic from inspect.go
}

func init() {
	rootCmd.AddCommand(inspectDirectCmd)
}
