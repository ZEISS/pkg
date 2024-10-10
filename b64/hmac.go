package b64

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// ContentHash computes the base64 encoded SHA-256 hash of the given content.
func ContentHash(content []byte) (string, error) {
	sha256 := sha256.New()

	_, err := sha256.Write(content)
	if err != nil {
		return "", err
	}

	hb := sha256.Sum(nil)
	s := base64.StdEncoding.EncodeToString(hb)

	return s, nil
}

// Hmac256 computes the base64 encoded HMAC-SHA-256 hash of the given message using the given secret.
func Hmac256(message, secret string) (string, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	hash := hmac.New(sha256.New, decodedSecret)

	_, err = hash.Write([]byte(message))
	if err != nil {
		return "", err
	}

	hb := hash.Sum(nil)
	s := base64.StdEncoding.EncodeToString(hb)

	return s, nil
}
