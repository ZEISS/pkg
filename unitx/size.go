package unitx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Integerer represents an integer.
type Integerer interface {
	// ToInt converts the integer to an integer.
	ToInt() (int, error)
	// ToInt64 converts the integer to an 64 bit integer.
	ToInt64() (int64, error)
	// ToInt32 converts the integer to an 32 bit integer.
	ToInt32() (int32, error)
}

// See: http://en.wikipedia.org/wiki/Binary_prefix
const (
	// Decimal.
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB

	// Binary.
	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
	PiB = 1024 * TiB
)

// HumanByteSize represents a human-readable byte size.
type HumanSize string

var _ Integerer = (*HumanSize)(nil)

// ToInt64 converts the human size to an 64 bit integer.
// nolint: gocyclo
func (h HumanSize) ToInt64() (int64, error) {
	var size int64

	lastDigit := 0
	for _, r := range h {
		if !unicode.IsDigit(r) && r != '-' {
			break
		}
		lastDigit++
	}

	numStr := string(h[:lastDigit])
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		var cerr *strconv.NumError
		if errors.As(err, &cerr) && errors.Is(cerr.Err, strconv.ErrSyntax) {
			return num, fmt.Errorf("integer '%s' is out of the range", h)
		}

		return num, fmt.Errorf("expected integer, but got '%s'", h)
	}

	suffix := strings.ToLower(strings.TrimSpace(string(h[lastDigit:])))

	switch suffix {
	case "":
		size = num * 1
	case "k":
		size = num * KB
	case "kb", "ki", "kib":
		size = num * KiB
	case "m":
		size = num * MB
	case "mb", "mi", "mib":
		size = num * MiB
	case "g":
		size = num * GB
	case "gb", "gi", "gib":
		size = num * GiB
	case "t":
		size = num * TB
	case "tb", "ti", "tib":
		size = num * TiB
	case "p":
		size = num * PB
	case "pb", "pi", "pib":
		size = num * PiB
	case "e":
		size = num * PB
	case "eb", "ei", "eib":
		size = num * PiB
	default:
		return 0, fmt.Errorf("expected valid unit, but got '%s'", h)
	}

	return size, nil
}

// ToInt converts the human size to an integer.
func (h HumanSize) ToInt() (int, error) {
	i64, err := h.ToInt64()
	if err != nil {
		return 0, err
	}

	return int(i64), nil
}

// ToInt32 converts the human size to an 32 bit integer.
func (h HumanSize) ToInt32() (int32, error) {
	i, err := h.ToInt64()
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}
