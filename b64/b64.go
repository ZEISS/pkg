package b64

import "encoding/base64"

// Base64 returns the base64 encoding of the input value.
func Base64[T string | []byte](value T) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
