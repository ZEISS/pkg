package errorx

import "k8s.io/apimachinery/pkg/api/errors"

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	return err != nil && errors.IsNotFound(err)
}

// IsAlreadyExists checks if an error is an already exists error.
func IsAlreadyExists(err error) bool {
	return err != nil && errors.IsAlreadyExists(err)
}

// IsConflict checks if an error is a conflict error.
func IsConflict(err error) bool {
	return err != nil && errors.IsConflict(err)
}

// IsInvalid checks if an error is an invalid error.
func IsInvalid(err error) bool {
	return err != nil && errors.IsInvalid(err)
}

// IsUnauthorized checks if an error is an unauthorized error.
func IsUnauthorized(err error) bool {
	return err != nil && errors.IsUnauthorized(err)
}

// IsForbidden checks if an error is a forbidden error.
func IsForbidden(err error) bool {
	return err != nil && errors.IsForbidden(err)
}

// IsServiceUnavailable checks if an error is a service unavailable error.
func IsServiceUnavailable(err error) bool {
	return err != nil && errors.IsServiceUnavailable(err)
}

// IsTimeout checks if an error is a timeout error.
func IsTimeout(err error) bool {
	return err != nil && errors.IsTimeout(err)
}

// IsServerTimeout checks if an error is a server timeout error.
func IsServerTimeout(err error) bool {
	return err != nil && errors.IsServerTimeout(err)
}

// IsGone checks if an error is a gone error.
func IsGone(err error) bool {
	return err != nil && errors.IsGone(err)
}

// IsBadRequest checks if an error is a bad request error.
func IsBadRequest(err error) bool {
	return err != nil && errors.IsBadRequest(err)
}
