package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Initializes OdinExePath
func init() {
	exePath, err := os.Executable()
	Must(err)
	OdinExePath = filepath.Dir(exePath)
}

// Gets a relative path from the Odin executable
func OdinPath(paths ...string) string {
	targetPath := filepath.Join(append([]string{OdinExePath}, paths...)...)
	return targetPath
}

// Helper func to check if file/dir Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	fmt.Printf("Could not check if %s exists", path)
	return false, err
}
