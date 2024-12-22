package filex

import "os"

// MkdirAll makes a new directory at path with the file mode.
func MkdirAll(path string, mode os.FileMode) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return err
	}

	err := os.MkdirAll(path, mode)
	if err != nil {
		return err
	}

	return nil
}
