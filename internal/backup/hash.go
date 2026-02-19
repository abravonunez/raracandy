package backup

import (
	"fmt"
	"os"
)

// CreateBackupWithHash creates a backup and saves its SHA256 hash
func CreateBackupWithHash(path string, hash string) error {
	// Create the backup
	if err := CreateBackup(path); err != nil {
		return err
	}

	// Save the hash
	hashPath := path + ".bak.sha256"
	hashContent := fmt.Sprintf("%s  %s.bak\n", hash, path)

	if err := os.WriteFile(hashPath, []byte(hashContent), 0644); err != nil {
		return fmt.Errorf("failed to write hash file: %w", err)
	}

	return nil
}

// VerifyBackupHash checks if a backup matches its saved hash
func VerifyBackupHash(path string, currentHash string) (bool, error) {
	hashPath := path + ".bak.sha256"

	data, err := os.ReadFile(hashPath)
	if err != nil {
		return false, fmt.Errorf("failed to read hash file: %w", err)
	}

	// Parse hash from file (format: "hash  filename")
	var savedHash string
	fmt.Sscanf(string(data), "%s", &savedHash)

	return savedHash == currentHash, nil
}
