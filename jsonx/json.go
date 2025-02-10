package jsonx

import (
	"encoding/json"

	"github.com/zeiss/pkg/errorx"
)

// Bytes is a type that represents a byte slice.
func Bytes(value any) []byte {
	return errorx.Must(json.Marshal(value))
}

// Prettify returns a pretty-printed JSON string.
func Prettify(value any) ([]byte, error) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}

	return json, nil
}
