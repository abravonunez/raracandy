package main

import "github.com/spf13/cobra"

var addItemDirectCmd = &cobra.Command{
	Use:   "add-item <save-file>",
	Short: "Add or modify an item in the bag (auto-detects version)",
	Long: `Add an item to the bag or update its quantity if it already exists.
Automatically detects whether the save file is from Pokémon Red, Blue, or Yellow.

If the item exists, its quantity will be updated. If it doesn't exist and there's
space in the bag, it will be added.

The save file will not be modified unless --out is specified.
Use --dry-run to preview changes without writing.`,
	Args: cobra.ExactArgs(1),
	RunE: runAddItem, // Reuse the same logic from add_item.go
}

func init() {
	rootCmd.AddCommand(addItemDirectCmd)

	addItemDirectCmd.Flags().StringVarP(&addItemOutput, "out", "o", "", "Output file path (required)")
	addItemDirectCmd.Flags().BoolVar(&addItemDryRun, "dry-run", false, "Preview changes without writing")
	addItemDirectCmd.Flags().StringVar(&addItemName, "item", "", "Item name (e.g., rare_candy)")
	addItemDirectCmd.Flags().IntVar(&addItemQty, "qty", 99, "Item quantity (1-99)")
	addItemDirectCmd.Flags().BoolVar(&addItemForce, "force", false, "Skip confirmation prompt")

	addItemDirectCmd.MarkFlagRequired("out")
	addItemDirectCmd.MarkFlagRequired("item")
}
