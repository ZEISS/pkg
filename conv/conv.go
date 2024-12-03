package conv

import (
	"fmt"
	"reflect"
	"strconv"
)

// String returns the string representation of the value.
func String(val any) string {
	switch val := val.(type) {
	case bool:
		return strconv.FormatBool(val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case []byte:
		return string(val)
	case string:
		return val
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Strings returns the string slice representation of the values.
func Strings[T any](vals ...T) []string {
	strs := make([]string, len(vals))
	for i, val := range vals {
		strs[i] = String(val)
	}

	return strs
}

// Bool returns the boolean representation of the value.
// nolint:gocyclo
func Bool(val any) bool {
	switch val := val.(type) {
	case bool:
		return val
	case *bool:
		if val == nil {
			return false
		}
		// Return actual value
		return *val
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return val != 0
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		if val == nil {
			return false
		}
		// Return integer boolean representation
		return reflect.ValueOf(val).Int() != 0
	case string:
		return val != ""
	case *string:
		if val == nil {
			return false
		}
		// Return string boolean representation
		return reflect.ValueOf(val).String() != ""
	case nil:
		return false
	default:
		panic("unknown type for Bool()")
	}
}

// Bools returns the boolean slice representation of the values.
func Bools[T any](vals ...T) []bool {
	bools := make([]bool, len(vals))
	for i, val := range vals {
		bools[i] = Bool(val)
	}

	return bools
}

// Int returns the integer representation of the value.
// nolint:gocyclo
func Int(val any) int {
	switch val := val.(type) {
	case bool:
		if val {
			return 1
		}

		return 0
	case int:
		return val
	case int8:
		return int(val)
	case int16:
		return int(val)
	case int32:
		return int(val)
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint8:
		return int(val)
	case uint16:
		return int(val)
	case uint32:
		return int(val)
	case uint64:
		return int(val)
	case float32:
		return int(val)
	case float64:
		return int(val)
	case *int:
		if val == nil {
			return 0
		}

		return *val
	case *int8:
		if val == nil {
			return 0
		}

		return int(*val)
	case *int16:
		if val == nil {
			return 0
		}

		return int(*val)
	case *int32:
		if val == nil {
			return 0
		}

		return int(*val)
	case *int64:
		if val == nil {
			return 0
		}

		return int(*val)
	case *uint:
		if val == nil {
			return 0
		}

		return int(*val)
	case *uint8:
		if val == nil {
			return 0
		}

		return int(*val)
	case *uint16:
		if val == nil {
			return 0
		}

		return int(*val)
	case *uint32:
		if val == nil {
			return 0
		}

		return int(*val)
	case *uint64:
		if val == nil {
			return 0
		}

		return int(*val)

	case *float32:
		if val == nil {
			return 0
		}

		return int(*val)
	case *float64:
		if val == nil {
			return 0
		}

		return int(*val)
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}

		return i
	case *string:
		if val == nil {
			return 0
		}

		i, err := strconv.Atoi(*val)
		if err != nil {
			panic(err)
		}

		return i
	case nil:
		return 0
	default:
		panic("unknown type for Int()")
	}
}

// Ints returns the integer slice representation of the values.
func Ints[T any](vals ...T) []int {
	ints := make([]int, len(vals))
	for i, val := range vals {
		ints[i] = Int(val)
	}

	return ints
}

// Bytes returns the byte slice representation of the value.
func Bytes(val any) []byte {
	switch val := val.(type) {
	case []byte:
		return val
	case string:
		return []byte(val)
	default:
		return []byte(String(val))
	}
}
