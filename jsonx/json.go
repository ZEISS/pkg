package jsonx

import (
	"encoding/json"

	"github.com/zeiss/pkg/errorx"
)

// Bytes is a type that represents a byte slice.
func Bytes(value any) []byte {
	return errorx.Must(json.Marshal(value))
}
