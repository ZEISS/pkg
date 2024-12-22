package filex

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// IsDirEmpty is a function to check if a directory is empty.
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(filepath.Clean(name))
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	_, err = f.Readdir(1)

	if errors.Is(err, io.EOF) {
		return true, nil
	}

	return false, err
}
