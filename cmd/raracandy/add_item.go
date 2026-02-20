package main

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/backup"
	"github.com/abravonunez/raracandy/internal/gen1/items"
	"github.com/abravonunez/raracandy/internal/gen1/save"
	"github.com/spf13/cobra"
)

var (
	addItemOutput string
	addItemDryRun bool
	addItemName   string
	addItemQty    int
	addItemForce  bool
)

var addItemCmd = &cobra.Command{
	Use:   "add-item <save-file>",
	Short: "Add or modify an item in the bag",
	Long: `Add an item to the bag or update its quantity if it already exists.
If the item exists, its quantity will be updated. If it doesn't exist and there's
space in the bag, it will be added.

The save file will not be modified unless --out is specified.
Use --dry-run to preview changes without writing.`,
	Args: cobra.ExactArgs(1),
	RunE: runAddItem,
}

func init() {
	yellowCmd.AddCommand(addItemCmd)

	addItemCmd.Flags().StringVarP(&addItemOutput, "out", "o", "", "Output file path (required)")
	addItemCmd.Flags().BoolVar(&addItemDryRun, "dry-run", false, "Preview changes without writing")
	addItemCmd.Flags().StringVar(&addItemName, "item", "", "Item name (e.g., rare_candy)")
	addItemCmd.Flags().IntVar(&addItemQty, "qty", 99, "Item quantity (1-99)")
	addItemCmd.Flags().BoolVar(&addItemForce, "force", false, "Skip confirmation prompt")

	addItemCmd.MarkFlagRequired("out")
	addItemCmd.MarkFlagRequired("item")
}

func runAddItem(cmd *cobra.Command, args []string) error {
	savePath := args[0]

	// Validate quantity
	if addItemQty < 1 || addItemQty > 99 {
		return fmt.Errorf("quantity must be between 1 and 99")
	}

	// Get item ID
	itemID, err := items.GetItemID(addItemName)
	if err != nil {
		return fmt.Errorf("invalid item: %w", err)
	}

	// Load save file
	fmt.Println("⚙️  Loading save...")
	s, err := save.Load(savePath)
	if err != nil {
		return fmt.Errorf("failed to load save: %w", err)
	}

	// Perform integrity check
	fmt.Println("🔍 Running integrity check...")
	report := s.CheckIntegrity()

	if !report.IsValid {
		fmt.Println("\n❌ Save file integrity check failed:")
		for _, err := range report.Errors {
			fmt.Printf("  • %s\n", err)
		}
		return fmt.Errorf("cannot modify corrupted save file")
	}

	fmt.Printf("✓ Integrity check passed\n")
	fmt.Printf("✓ Detected: %s\n", report.GameVersion)
	fmt.Println()

	// Find current state
	currentIdx := items.FindItemIndex(s, itemID)
	var currentQty byte = 0
	if currentIdx >= 0 {
		bagItems := items.GetBagItems(s)
		currentQty = bagItems[currentIdx].Quantity
	}

	// Preview changes
	itemName := items.GetItemName(itemID)
	fmt.Println("Changes to be applied:")
	fmt.Println("  Bag items:")
	if currentIdx >= 0 {
		fmt.Printf("    - %s: %d → %d (%+d)\n", itemName, currentQty, addItemQty, addItemQty-int(currentQty))
	} else {
		fmt.Printf("    - %s: (new) → %d\n", itemName, addItemQty)
	}

	oldChecksum := s.GetChecksum()
	fmt.Printf("  Checksum: 0x%02X → (will recalculate)\n", oldChecksum)

	if addItemDryRun {
		fmt.Println("\n[DRY RUN] No changes written")
		return nil
	}

	// Ask for confirmation if not in force mode
	if !addItemForce {
		changes := []string{
			fmt.Sprintf("Add/modify %s to quantity %d", itemName, addItemQty),
			"Recalculate checksum",
		}
		if !save.ConfirmWithDetails(changes) {
			fmt.Println("\n❌ Operation cancelled by user")
			return nil
		}
	}

	// Apply changes
	fmt.Println("\n✍️  Applying changes...")
	if err := items.SetItemQuantity(s, itemID, byte(addItemQty)); err != nil {
		return fmt.Errorf("failed to set item: %w", err)
	}

	// Get original hash before backup
	originalHash := s.GetSHA256()

	// Create backup with hash
	fmt.Println("💾 Creating backup...")
	if err := backup.CreateBackupWithHash(savePath, originalHash); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	fmt.Printf("✓ Backup created: %s\n", backup.GetBackupPath(savePath))
	fmt.Printf("✓ Backup hash saved: %s.bak.sha256\n", savePath)

	// Write output
	if err := s.Write(addItemOutput); err != nil {
		return fmt.Errorf("failed to write save: %w", err)
	}

	// Verify written file
	written, err := save.Load(addItemOutput)
	if err != nil {
		return fmt.Errorf("failed to verify written file: %w", err)
	}
	if !written.ValidateChecksum() {
		return fmt.Errorf("verification failed: checksum invalid after write")
	}

	newChecksum := s.GetChecksum()
	fmt.Printf("\n✓ Save written: %s\n", addItemOutput)
	fmt.Printf("✓ Checksum updated: 0x%02X → 0x%02X\n", oldChecksum, newChecksum)
	fmt.Printf("✓ Verification passed\n")
	fmt.Printf("\n🎉 Success! Your save is ready to use.")

	return nil
}
