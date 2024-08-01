package errorx

import (
	"errors"
	"fmt"
)

// RecoverError converts a panic value to an error.
func RecoverError(r interface{}) error {
	switch x := r.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return errors.New(fmt.Sprint(x))
	}
}
