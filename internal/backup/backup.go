package backup

import (
	"fmt"
	"os"
	"time"
)

// CreateBackup creates a backup of the specified file
// The backup is created with a .bak extension
func CreateBackup(path string) error {
	backupPath := path + ".bak"
	return createBackupWithPath(path, backupPath)
}

// CreateTimestampedBackup creates a backup with a timestamp
// Format: original.sav.bak.20060102-150405
func CreateTimestampedBackup(path string) error {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.bak.%s", path, timestamp)
	return createBackupWithPath(path, backupPath)
}

// createBackupWithPath performs the actual backup operation
func createBackupWithPath(srcPath, dstPath string) error {
	// Read source file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write backup file
	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// BackupExists checks if a backup file exists
func BackupExists(path string) bool {
	backupPath := path + ".bak"
	_, err := os.Stat(backupPath)
	return err == nil
}

// GetBackupPath returns the path for a backup file
func GetBackupPath(path string) string {
	return path + ".bak"
}
