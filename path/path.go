package path

import (
	"fmt"
	"os"
	"path/filepath"
)

var OdinExePath string

// Initializes OdinExePath. All paths used by Odin will be relative to this.
func init() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	OdinExePath = filepath.Dir(exePath)
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

// Gets a relative path from the Odin executable
func Relative(paths ...string) string {
	targetPath := filepath.Join(append([]string{OdinExePath}, paths...)...)
	return targetPath
}
