package main

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/gen1/items"
	"github.com/abravonunez/raracandy/internal/gen1/money"
	"github.com/abravonunez/raracandy/internal/gen1/save"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <save-file>",
	Short: "Inspect a Pokémon Yellow save file",
	Long:  `Display information about a Pokémon Yellow save file including money, bag items, and checksum status.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInspect,
}

func init() {
	yellowCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	savePath := args[0]

	// Load save file
	s, err := save.Load(savePath)
	if err != nil {
		return fmt.Errorf("failed to load save: %w", err)
	}

	fmt.Printf("Save File: %s\n", savePath)
	fmt.Printf("Size: %d KB\n", len(s.Data())/1024)
	fmt.Println()

	// Checksum info
	stored := s.GetChecksum()
	calculated := s.CalculateChecksum()
	checksumValid := s.ValidateChecksum()

	fmt.Println("Checksum:")
	fmt.Printf("  Stored:     0x%02X\n", stored)
	fmt.Printf("  Calculated: 0x%02X\n", calculated)
	if checksumValid {
		fmt.Println("  Status:     ✓ Valid")
	} else {
		fmt.Println("  Status:     ✗ Invalid (file may be corrupted)")
	}
	fmt.Println()

	// Money
	playerMoney := money.GetMoney(s)
	fmt.Printf("Money: %s\n", money.FormatMoney(playerMoney))
	fmt.Println()

	// Bag items
	bagItems := items.GetBagItems(s)
	fmt.Printf("Bag (%d/%d items):\n", len(bagItems), items.MaxBagItems)
	if len(bagItems) == 0 {
		fmt.Println("  (empty)")
	} else {
		for _, item := range bagItems {
			fmt.Printf("  - %s x%d\n", item.Name, item.Quantity)
		}
	}

	return nil
}
