package filex

import (
	"os"
)

// TearDownFunc ...
type TearDownFunc func()

// TempDir ...
func TempDir() (*os.File, TearDownFunc, error) {
	f, err := os.CreateTemp(os.TempDir(), "streams_")
	if err != nil {
		return f, nil, err
	}

	fn := func() {
		_ = os.Remove(f.Name())
	}

	return f, fn, nil
}
