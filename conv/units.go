package conv

// ByteSize represents the size of a value in bits.
type ByteSize int64

const (
	B  ByteSize = 1
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

// ByteSizeUnits is a map of byte sizes.
var ByteSizeUnits = map[string]ByteSize{
	"B":  B,
	"KB": KB,
	"MB": MB,
	"GB": GB,
	"TB": TB,
	"PB": PB,
	"EB": EB,
}

// ByteSizes converts a value of float64 to a full integer value of a byte
// based on a unit size.
func ByteSizes(b float64, size string) int64 {
	s, ok := ByteSizeUnits[size]
	if !ok {
		return 0
	}

	return int64(b * float64(s))
}
