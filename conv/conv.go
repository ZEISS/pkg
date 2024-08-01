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
	case string:
		return val
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Bool returns the boolean representation of the value.
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
