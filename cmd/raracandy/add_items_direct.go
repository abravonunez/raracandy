package main

import "github.com/spf13/cobra"

var addItemsDirectCmd = &cobra.Command{
	Use:   "add-items <save-file>",
	Short: "Add multiple items to the bag in a single operation (auto-detects version)",
	Long: `Add multiple items to the bag or update their quantities if they already exist.
Automatically detects whether the save file is from Pokémon Red, Blue, or Yellow.

This is an atomic operation - all items are added in a single transaction.
If any item fails to be added, the entire operation is rolled back.

The save file will not be modified unless --out is specified.
Use --dry-run to preview changes without writing.`,
	Args: cobra.ExactArgs(1),
	RunE: runAddItems, // Reuse the same logic from add_items.go
}

func init() {
	rootCmd.AddCommand(addItemsDirectCmd)

	addItemsDirectCmd.Flags().StringVarP(&addItemsOutput, "out", "o", "", "Output file path (required)")
	addItemsDirectCmd.Flags().BoolVar(&addItemsDryRun, "dry-run", false, "Preview changes without writing")
	addItemsDirectCmd.Flags().StringSliceVar(&addItemsNames, "item", []string{}, "Item names (can be specified multiple times)")
	addItemsDirectCmd.Flags().IntSliceVar(&addItemsQtys, "qty", []int{}, "Item quantities (must match number of items)")
	addItemsDirectCmd.Flags().BoolVar(&addItemsForce, "force", false, "Skip confirmation prompt")

	addItemsDirectCmd.MarkFlagRequired("out")
	addItemsDirectCmd.MarkFlagRequired("item")
	addItemsDirectCmd.MarkFlagRequired("qty")
}
