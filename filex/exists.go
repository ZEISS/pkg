package filex

import (
	"os"
)

// FileExists is testing if a file exists at a given path.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}

	return true, nil
}

// FileNotExists is testing if a file does not exists at a given path.
func FileNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if err == nil {
		return false, nil
	}

	return true, nil
}
