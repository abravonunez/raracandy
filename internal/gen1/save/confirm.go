package save

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmOperation asks the user to confirm a dangerous operation
func ConfirmOperation(message string) bool {
	fmt.Printf("\n⚠️  WARNING: %s\n", message)
	fmt.Print("Type 'yes' to continue: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes"
}

// ConfirmWithDetails shows detailed changes and asks for confirmation
func ConfirmWithDetails(changes []string) bool {
	fmt.Println("\n📝 The following changes will be made:")
	for _, change := range changes {
		fmt.Printf("  • %s\n", change)
	}
	fmt.Println()

	return ConfirmOperation("You are about to modify your save file")
}
