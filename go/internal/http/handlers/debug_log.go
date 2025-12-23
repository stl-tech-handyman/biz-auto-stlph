package handlers

import (
	"os"
	"path/filepath"
)

func GetLogPath() string {
	// Try to find .cursor/debug.log relative to workspace root
	// Server runs from go/ directory, so go up one level
	if wd, err := os.Getwd(); err == nil {
		// If we're in go/, go up to workspace root
		if filepath.Base(wd) == "go" {
			return filepath.Join(wd, "..", ".cursor", "debug.log")
		}
		// Otherwise assume we're in workspace root
		return filepath.Join(wd, ".cursor", "debug.log")
	}
	// Fallback
	return "../../.cursor/debug.log"
}

