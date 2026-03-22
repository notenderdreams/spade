package utils

import (
	"io"
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

func GetScriptsDir() string {
	dir := getSpadeConfigDir()
	scriptsDir := filepath.Join(dir, "scripts")
	_ = os.MkdirAll(scriptsDir, 0755)
	return scriptsDir
}

func GetScriptPath(name, path string) string {
	return filepath.Join(GetScriptsDir(), name+"_"+filepath.Base(path))
}

func StoreAtScriptDir(name, srcPath string) error {
	dstPath := GetScriptPath(name, srcPath)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func ScriptExistsOrBackup(name, path string) (string, error) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return path, nil
	}

	backupPath := GetScriptPath(name, path)
	if info, err := os.Stat(backupPath); err == nil && !info.IsDir() {
		return backupPath, nil
	}

	return "", os.ErrNotExist
}