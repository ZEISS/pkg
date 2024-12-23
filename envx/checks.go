package envx

import (
	"context"
	"fmt"
	"os/user"

	"github.com/zeiss/pkg/filex"
)

// HasUser is checking if the current user is available.
func HasUser() Check {
	return func(ctx context.Context) error {
		_, err := user.Current()
		if err != nil {
			return err
		}

		return nil
	}
}

// IsDirEmpty is checking if the folder is not empty.
func IsDirEmpty(path string) Check {
	return func(ctx context.Context) error {
		empty, err := filex.IsDirEmpty(path)
		if err != nil {
			return err
		}

		if !empty {
			return fmt.Errorf("folder %s is not empty", path)
		}

		return nil
	}
}
