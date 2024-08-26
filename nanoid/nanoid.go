package nanoid

import (
	crand "crypto/rand"
	"errors"
	"math/bits"
	"sync"
)

var ErrInvalidLength = errors.New("nanoid: length for ID is invalid (must be within 2-255)")

type generator = func() string

// `A-Za-z0-9_-`.
var defaultAlphabet = [64]byte{
	'A', 'B', 'C', 'D', 'E',
	'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y',
	'Z', 'a', 'b', 'c', 'd',
	'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's',
	't', 'u', 'v', 'w', 'x',
	'y', 'z', '0', '1', '2',
	'3', '4', '5', '6', '7',
	'8', '9', '-', '_',
}

// DefaultLength is the default length for NanoID.
const DefaultLength = 21

// New returns a new NanoID generator with the standard alphabet and length.
func New(length int) (generator, error) {
	if length < 2 || length > 255 {
		return nil, ErrInvalidLength
	}

	size := length * length * 7
	b := make([]byte, size)
	crand.Read(b)

	offset := 0
	id := make([]byte, length)

	var mu sync.Mutex

	return func() string {
		mu.Lock()
		defer mu.Unlock()

		if offset == size {
			crand.Read(b)
			offset = 0
		}

		for i := 0; i < length; i++ {
			id[i] = defaultAlphabet[b[i+offset]&63]
		}

		offset += length

		return string(id)
	}, nil
}

// Unicode returns a Nano ID generator which uses a custom Unicode alphabet.
func Unicode(alphabet string, length int) (generator, error) {
	if length < 2 || length > 255 {
		return nil, ErrInvalidLength
	}

	alphabetLen := len(alphabet)
	runes := []rune(alphabet)

	x := uint32(alphabetLen) - 1
	clz := bits.LeadingZeros32(x | 1)
	mask := (2 << (31 - clz)) - 1
	step := (length / 5) * 8

	b := make([]byte, step)
	id := make([]rune, length)

	j, idx := 0, 0

	var mu sync.Mutex

	return func() string {
		mu.Lock()
		defer mu.Unlock()

		for {
			crand.Read(b)
			for i := 0; i < step; i++ {
				idx = int(b[i]) & mask
				if idx < alphabetLen {
					id[j] = runes[idx]
					j++
					if j == length {
						j = 0
						return string(id)
					}
				}
			}
		}
	}, nil
}
