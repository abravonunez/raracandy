package main

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/backup"
	"github.com/abravonunez/raracandy/internal/gen1/items"
	"github.com/abravonunez/raracandy/internal/gen1/save"
	"github.com/spf13/cobra"
)

var (
	addItemsOutput  string
	addItemsDryRun  bool
	addItemsNames   []string
	addItemsQtys    []int
	addItemsForce   bool
)

var addItemsCmd = &cobra.Command{
	Use:   "add-items <save-file>",
	Short: "Add or modify multiple items in the bag",
	Long: `Add multiple items to the bag or update their quantities if they already exist.
You can specify multiple --item and --qty flags to add multiple items in a single operation.

Examples:
  # Add multiple items with different quantities
  raracandy yellow add-items save.sav -o output.sav \
    --item rare_candy --qty 99 \
    --item master_ball --qty 50 \
    --item ultra_ball --qty 80

  # Preview changes without writing
  raracandy yellow add-items save.sav -o output.sav \
    --item rare_candy --qty 99 \
    --item potion --qty 50 \
    --dry-run

The save file will not be modified unless --out is specified.
Use --dry-run to preview changes without writing.`,
	Args: cobra.ExactArgs(1),
	RunE: runAddItems,
}

func init() {
	yellowCmd.AddCommand(addItemsCmd)

	addItemsCmd.Flags().StringVarP(&addItemsOutput, "out", "o", "", "Output file path (required)")
	addItemsCmd.Flags().BoolVar(&addItemsDryRun, "dry-run", false, "Preview changes without writing")
	addItemsCmd.Flags().StringSliceVar(&addItemsNames, "item", []string{}, "Item name (can be repeated)")
	addItemsCmd.Flags().IntSliceVar(&addItemsQtys, "qty", []int{}, "Item quantity 1-99 (can be repeated, must match number of items)")
	addItemsCmd.Flags().BoolVar(&addItemsForce, "force", false, "Skip confirmation prompt")

	addItemsCmd.MarkFlagRequired("out")
	addItemsCmd.MarkFlagRequired("item")
}

type itemChange struct {
	name       string
	itemID     byte
	newQty     int
	currentQty byte
	isNew      bool
}

func runAddItems(cmd *cobra.Command, args []string) error {
	savePath := args[0]

	// Validate that we have items to add
	if len(addItemsNames) == 0 {
		return fmt.Errorf("at least one --item must be specified")
	}

	// Validate that quantities match items
	if len(addItemsQtys) != len(addItemsNames) {
		return fmt.Errorf("number of --qty flags (%d) must match number of --item flags (%d)",
			len(addItemsQtys), len(addItemsNames))
	}

	// Validate all quantities
	for i, qty := range addItemsQtys {
		if qty < 1 || qty > 99 {
			return fmt.Errorf("quantity for item %d (%s) must be between 1 and 99, got %d",
				i+1, addItemsNames[i], qty)
		}
	}

	// Build list of changes and validate all items exist
	changes := make([]itemChange, 0, len(addItemsNames))
	for i, itemName := range addItemsNames {
		itemID, err := items.GetItemID(itemName)
		if err != nil {
			return fmt.Errorf("invalid item %d: %w", i+1, err)
		}

		changes = append(changes, itemChange{
			name:   items.GetItemName(itemID),
			itemID: itemID,
			newQty: addItemsQtys[i],
		})
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

	// Find current state for all items
	for i := range changes {
		currentIdx := items.FindItemIndex(s, changes[i].itemID)
		if currentIdx >= 0 {
			bagItems := items.GetBagItems(s)
			changes[i].currentQty = bagItems[currentIdx].Quantity
			changes[i].isNew = false
		} else {
			changes[i].currentQty = 0
			changes[i].isNew = true
		}
	}

	// Preview changes
	fmt.Println("Changes to be applied:")
	fmt.Println("  Bag items:")
	for _, change := range changes {
		if change.isNew {
			fmt.Printf("    - %s: (new) → %d\n", change.name, change.newQty)
		} else {
			delta := change.newQty - int(change.currentQty)
			fmt.Printf("    - %s: %d → %d (%+d)\n",
				change.name, change.currentQty, change.newQty, delta)
		}
	}

	oldChecksum := s.GetChecksum()
	fmt.Printf("  Checksum: 0x%02X → (will recalculate)\n", oldChecksum)

	if addItemsDryRun {
		fmt.Println("\n[DRY RUN] No changes written")
		return nil
	}

	// Ask for confirmation if not in force mode
	if !addItemsForce {
		changeList := make([]string, len(changes)+1)
		for i, change := range changes {
			if change.isNew {
				changeList[i] = fmt.Sprintf("Add %s (qty: %d)", change.name, change.newQty)
			} else {
				changeList[i] = fmt.Sprintf("Update %s to quantity %d", change.name, change.newQty)
			}
		}
		changeList[len(changes)] = "Recalculate checksum"

		if !save.ConfirmWithDetails(changeList) {
			fmt.Println("\n❌ Operation cancelled by user")
			return nil
		}
	}

	// Apply all changes
	fmt.Println("\n✍️  Applying changes...")
	for i, change := range changes {
		if err := items.SetItemQuantity(s, change.itemID, byte(change.newQty)); err != nil {
			return fmt.Errorf("failed to set item %d (%s): %w", i+1, change.name, err)
		}
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
	if err := s.Write(addItemsOutput); err != nil {
		return fmt.Errorf("failed to write save: %w", err)
	}

	// Verify written file
	written, err := save.Load(addItemsOutput)
	if err != nil {
		return fmt.Errorf("failed to verify written file: %w", err)
	}
	if !written.ValidateChecksum() {
		return fmt.Errorf("verification failed: checksum invalid after write")
	}

	newChecksum := s.GetChecksum()
	fmt.Printf("\n✓ Save written: %s\n", addItemsOutput)
	fmt.Printf("✓ Checksum updated: 0x%02X → 0x%02X\n", oldChecksum, newChecksum)
	fmt.Printf("✓ Verification passed\n")
	fmt.Printf("✓ %d item(s) added/updated\n", len(changes))
	fmt.Printf("\n🎉 Success! Your save is ready to use.")

	return nil
}
