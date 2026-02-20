package main

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/backup"
	"github.com/abravonunez/raracandy/internal/gen1/money"
	"github.com/abravonunez/raracandy/internal/gen1/save"
	"github.com/spf13/cobra"
)

var (
	setMoneyOutput string
	setMoneyDryRun bool
	setMoneyAmount int
	setMoneyForce  bool
)

var setMoneyCmd = &cobra.Command{
	Use:   "set-money <save-file>",
	Short: "Set the player's money",
	Long: `Set the player's money to a specific amount.
The maximum amount is 999,999.

The save file will not be modified unless --out is specified.
Use --dry-run to preview changes without writing.`,
	Args: cobra.ExactArgs(1),
	RunE: runSetMoney,
}

func init() {
	yellowCmd.AddCommand(setMoneyCmd)

	setMoneyCmd.Flags().StringVarP(&setMoneyOutput, "out", "o", "", "Output file path (required)")
	setMoneyCmd.Flags().BoolVar(&setMoneyDryRun, "dry-run", false, "Preview changes without writing")
	setMoneyCmd.Flags().IntVar(&setMoneyAmount, "amount", 0, "Money amount (0-999999)")
	setMoneyCmd.Flags().BoolVar(&setMoneyForce, "force", false, "Skip confirmation prompt")

	setMoneyCmd.MarkFlagRequired("out")
	setMoneyCmd.MarkFlagRequired("amount")
}

func runSetMoney(cmd *cobra.Command, args []string) error {
	savePath := args[0]

	// Validate amount
	if setMoneyAmount < 0 || setMoneyAmount > money.MaxMoney {
		return fmt.Errorf("amount must be between 0 and %d", money.MaxMoney)
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

	// Get current money
	currentMoney := money.GetMoney(s)

	// Preview changes
	fmt.Println("Changes to be applied:")
	fmt.Printf("  Money: %s → %s (%+d)\n",
		money.FormatMoney(currentMoney),
		money.FormatMoney(uint32(setMoneyAmount)),
		setMoneyAmount-int(currentMoney))

	oldChecksum := s.GetChecksum()
	fmt.Printf("  Checksum: 0x%02X → (will recalculate)\n", oldChecksum)

	if setMoneyDryRun {
		fmt.Println("\n[DRY RUN] No changes written")
		return nil
	}

	// Ask for confirmation if not in force mode
	if !setMoneyForce {
		changes := []string{
			fmt.Sprintf("Set money to %s", money.FormatMoney(uint32(setMoneyAmount))),
			"Recalculate checksum",
		}
		if !save.ConfirmWithDetails(changes) {
			fmt.Println("\n❌ Operation cancelled by user")
			return nil
		}
	}

	// Apply changes
	fmt.Println("\n✍️  Applying changes...")
	if err := money.SetMoney(s, uint32(setMoneyAmount)); err != nil {
		return fmt.Errorf("failed to set money: %w", err)
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
	if err := s.Write(setMoneyOutput); err != nil {
		return fmt.Errorf("failed to write save: %w", err)
	}

	// Verify written file
	written, err := save.Load(setMoneyOutput)
	if err != nil {
		return fmt.Errorf("failed to verify written file: %w", err)
	}
	if !written.ValidateChecksum() {
		return fmt.Errorf("verification failed: checksum invalid after write")
	}

	newChecksum := s.GetChecksum()
	fmt.Printf("\n✓ Save written: %s\n", setMoneyOutput)
	fmt.Printf("✓ Checksum updated: 0x%02X → 0x%02X\n", oldChecksum, newChecksum)
	fmt.Printf("✓ Verification passed\n")
	fmt.Printf("\n🎉 Success! Your save is ready to use.")

	return nil
}
