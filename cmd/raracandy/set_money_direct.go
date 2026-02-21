package main

import "github.com/spf13/cobra"

var setMoneyDirectCmd = &cobra.Command{
	Use:   "set-money <save-file>",
	Short: "Set player money (auto-detects version)",
	Long: `Set the player's money to a specific amount.
Automatically detects whether the save file is from Pokémon Red, Blue, or Yellow.

Money is stored in BCD (Binary-Coded Decimal) format and can range from 0 to 999,999.

The save file will not be modified unless --out is specified.
Use --dry-run to preview changes without writing.`,
	Args: cobra.ExactArgs(1),
	RunE: runSetMoney, // Reuse the same logic from set_money.go
}

func init() {
	rootCmd.AddCommand(setMoneyDirectCmd)

	setMoneyDirectCmd.Flags().StringVarP(&setMoneyOutput, "out", "o", "", "Output file path (required)")
	setMoneyDirectCmd.Flags().BoolVar(&setMoneyDryRun, "dry-run", false, "Preview changes without writing")
	setMoneyDirectCmd.Flags().IntVar(&setMoneyAmount, "amount", 0, "Amount of money (0-999999)")
	setMoneyDirectCmd.Flags().BoolVar(&setMoneyForce, "force", false, "Skip confirmation prompt")

	setMoneyDirectCmd.MarkFlagRequired("out")
	setMoneyDirectCmd.MarkFlagRequired("amount")
}
