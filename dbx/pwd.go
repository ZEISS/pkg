package dbx

import (
	"errors"

	"github.com/zeiss/pkg/conv"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrFailedToHashPassword is returned when the password hashing fails
	ErrFailedToHashPassword = errors.New("dbx: failed to hash password")
	// ErrFailedCheckPassword is returned when the password check fails
	ErrFailedCheckPassword = errors.New("dbx: failed to check password")
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Join(ErrFailedToHashPassword, err)
	}

	return conv.String(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.Join(ErrFailedCheckPassword, err)
	}

	return nil
}
