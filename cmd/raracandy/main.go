package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "raracandy",
	Short: "A CLI tool to safely edit Pokémon Gen 1 save files",
	Long: `raracandy is a command-line tool for editing Pokémon Yellow save files.
It allows you to modify items, money, and other game data while maintaining
save file integrity through proper checksum calculation.

Never distributes or modifies ROMs - only operates on save files you own.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
