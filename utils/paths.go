package utils

import (
	"os"
	"path/filepath"
)

func getSpadeConfigDir() string {
	dir, _ := os.UserConfigDir()
	spadeDir := filepath.Join(dir, "spade")
	_ = os.MkdirAll(spadeDir, 0755)
	return spadeDir
}

func GetDBPath() string {
	return filepath.Join(getSpadeConfigDir(), "spade.db")
}
