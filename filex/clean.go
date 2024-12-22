package filex

import "os"

// Clean cleans the path and recreates it with the given mode.
func Clean(path string, mode os.FileMode) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = MkdirAll(path, mode)
	if err != nil {
		return err
	}

	return nil
}
