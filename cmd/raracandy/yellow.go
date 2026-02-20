package main

import (
	"github.com/spf13/cobra"
)

var yellowCmd = &cobra.Command{
	Use:   "yellow",
	Short: "Commands for Pokémon Yellow save files",
	Long:  `Edit and inspect Pokémon Yellow (Gen 1) save files.`,
}

func init() {
	rootCmd.AddCommand(yellowCmd)
}
