package dbx

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrFailedToHashPassword is returned when the password hashing fails.
	ErrFailedToHashPassword = errors.New("dbx: failed to hash password")
	// ErrFailedCheckPassword is returned when the password check fails.
	ErrFailedCheckPassword = errors.New("dbx: failed to check password")
)

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Join(ErrFailedToHashPassword, err)
	}

	return hashedPassword, nil
}

// CheckPassword checks if the provided password is correct or not.
func CheckPassword(password []byte, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return errors.Join(ErrFailedCheckPassword, err)
	}

	return nil
}
