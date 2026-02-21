package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "raracandy",
	Short: "A CLI tool to safely edit Pokémon Gen 1 save files",
	Long: `raracandy is a command-line tool for editing Pokémon Red/Blue/Yellow save files.
It automatically detects your game version and applies the correct offsets.
You can modify items, money, and other game data while maintaining save file
integrity through proper checksum calculation.

Commands can be used directly (auto-detect version):
  raracandy add-item pokemon.sav --item rare_candy --qty 99 --out modified.sav

Or with explicit version (backward compatibility):
  raracandy yellow add-item pokemon.sav --item rare_candy --qty 99 --out modified.sav

Never distributes or modifies ROMs - only operates on save files you own.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
