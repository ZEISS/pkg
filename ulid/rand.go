package ulid

import (
	"errors"
	"io"
)

var (
	// ErrTimeOverflow is returned if a timestamp >= max time is provided to MonotonicReader.
	ErrTimeOverflow = errors.New("ulid: timestamp overflows max time")

	// ErrNegativeTime is returned if a timestamp predating the previous is provided to MonotonicReader.
	ErrNegativeTime = errors.New("ulid: timestamp predates previous")

	// ErrMonotonicOverflow is returned when incrementing previous ULID's entropy bytes would result
	// in an entropy overflow. The solution is to call .Next() again with a later timestamp.
	ErrMonotonicOverflow = errors.New("ulid: monotonic entropy overflow")
)

// global monotonic reader.
var reader *MonotonicReader

// SetSource updates the global MonotonicReader's source. Note that the default
// MonotonicReader source is cryptographically safe (uses "crypto/rand.Reader").
func SetSource(src io.Reader) error {
	mutex.Lock()
	defer mutex.Unlock()

	err := setreader(src)
	if err != nil {
		return err
	}

	return nil
}

// setreader will set global MonotonicReader's source, ONLY on success.
func setreader(src io.Reader) (err error) {
	reader, err = NewMonotonicReader(src)
	if err != nil {
		reader = nil
	}
	return
}

// MonotonicReader wraps a reader to provide montonically increasing ULID values.
type MonotonicReader struct {
	src io.Reader  // entropy source reader
	buf [1024]byte // buffered entropy from reader
	idx int        // index in buffered entropy
	ent uint80     // current entropy value
	ovf int64      // timestamp with entropy overflow
	ts  int64      // timestamp of last entropy read
}

// NewMonotonicReader returns a new instance for given entropy source and read buffer
// size. Note that 'src' MUST yield random bytes or else monotonic reads are not guaranteed
// to terminate, since there may not be enough entropy to compute a monotonic increment.
func NewMonotonicReader(src io.Reader) (*MonotonicReader, error) {
	// Create new reader
	r := &MonotonicReader{
		ts:  Now(),
		src: src,
	} // set empty buffer
	r.idx = len(r.buf) - 1

	// Force an initial entropy buffer fill
	if err := r.refill(10); err != nil {
		return nil, err
	}

	// Copy direct to uint80
	copy(r.ent[:], r.buf[:10])
	r.idx = 10

	return r, nil
}

// NextReverse calculates next available ULID for given Timestamp.
func (r *MonotonicReader) NextReverse(ts int64, dst *ULID) error {
	switch {
	case ts >= maxTime():
		return ErrTimeOverflow

	case ts == r.ovf:
		return ErrMonotonicOverflow
	}

	ts = maxTime() - ts

	// Set ULID timestamp
	// (converts to bigendian)
	dst[0] = byte(ts >> 40)
	dst[1] = byte(ts >> 32)
	dst[2] = byte(ts >> 24)
	dst[3] = byte(ts >> 16)
	dst[4] = byte(ts >> 8)
	dst[5] = byte(ts)

	// Read entropy data into ULID
	return r.entropy(ts, dst)
}

// Next calculates next available ULID for given Timestamp.
func (r *MonotonicReader) Next(ts int64, dst *ULID) error {
	switch {
	case ts < r.ts:
		return ErrNegativeTime

	case ts >= maxTime():
		return ErrTimeOverflow

	case ts == r.ovf:
		return ErrMonotonicOverflow
	}

	// Set ULID timestamp
	// (converts to bigendian)
	dst[0] = byte(ts >> 40)
	dst[1] = byte(ts >> 32)
	dst[2] = byte(ts >> 24)
	dst[3] = byte(ts >> 16)
	dst[4] = byte(ts >> 8)
	dst[5] = byte(ts)

	// Read entropy data into ULID
	return r.entropy(ts, dst)
}

// entropy reads next monotonic entropy into last 10 bytes of ULID.
func (r *MonotonicReader) entropy(ts int64, dst *ULID) error {
	// ULID entropy bytes
	const bitlen = 10

	if ts == r.ts {
		// Get entropy increment
		inc, err := r.random()
		if err != nil {
			return err
		}

		// Increment prev entropy value
		if !r.ent.add(uint64(inc)) {
			r.ovf = ts // set overflow ts
			return ErrMonotonicOverflow
		}

		// Set return value
		r.ent.load(dst[6:])

		return nil
	}

	// Check if entropy buffer needs refill
	if err := r.refill(bitlen); err != nil {
		return err
	}

	// Copy direct to ULID
	next := r.idx + bitlen
	copy(dst[6:], r.buf[r.idx:next])
	r.idx = next

	// Unset most significant bit of ULID
	// entropy to ensure that even for the
	// largest entropy value, there's still
	// room to generate more ULIDs in this ms.
	dst[7] = 0

	// Set new entropy
	r.ent.store(dst[6:])

	// Update ts
	r.ts = ts

	return nil
}

// random returns a random value in range [1, ^uint32(0)].
func (r *MonotonicReader) random() (u uint32, err error) {
	// bits in uint32
	const bitlen = 4

	// Check if entropy buffer needs refill
	if err = r.refill(bitlen); err != nil {
		return 0, err
	}

	// Get slice into entropy buf
	next := r.idx + bitlen
	buf := r.buf[r.idx:next]
	r.idx = next

	// Convert random bytes into uint32
	// (converts from little-endian)
	u = 1 + uint32(buf[0]) |
		uint32(buf[1])<<8 |
		uint32(buf[2])<<16 |
		uint32(buf[3])<<24

	return
}

// refill will guarantee that at least 'require' bytes are left in entropy buffer.
func (r *MonotonicReader) refill(require int) error {
	// No refil is required for read
	if require < len(r.buf)-r.idx {
		return nil
	}

	// Perform full read into buffer
	_, err := io.ReadFull(r.src, r.buf[:])
	if err != nil {
		return err
	}

	// Reset idx
	r.idx = 0

	return nil
}

// uint80 is an unsigned integer of 80 bits.
type uint80 [10]uint8

// store will set value of uint80 into bytes.
func (u *uint80) store(src []byte) {
	_ = src[9] // we know len
	u[0] = src[0]
	u[1] = src[1]
	u[2] = src[2]
	u[3] = src[3]
	u[4] = src[4]
	u[5] = src[5]
	u[6] = src[6]
	u[7] = src[7]
	u[8] = src[8]
	u[9] = src[9]
}

// load will load value of uint80 into bytes.
func (u *uint80) load(dst []byte) {
	_ = dst[9] // we know len
	dst[0] = u[0]
	dst[1] = u[1]
	dst[2] = u[2]
	dst[3] = u[3]
	dst[4] = u[4]
	dst[5] = u[5]
	dst[6] = u[6]
	dst[7] = u[7]
	dst[8] = u[8]
	dst[9] = u[9]
}

// hi returns the high 16 bits of uint80.
func (u uint80) hi() uint16 {
	// converts from bigendian
	return uint16(u[1]) | uint16(u[0])<<8
}

// setHi sets the high 16 bits of uint80.
func (u *uint80) setHi(hi uint16) {
	// converts to bigendian
	u[0] = byte(hi >> 8)
	u[1] = byte(hi)
}

// lo returns the low 64 bits of uint80.
func (u uint80) lo() uint64 {
	// converts from bigendian
	return uint64(u[9]) | uint64(u[8])<<8 | uint64(u[7])<<16 | uint64(u[6])<<24 |
		uint64(u[5])<<32 | uint64(u[4])<<40 | uint64(u[3])<<48 | uint64(u[2])<<56
}

// setLo sets the low 64 bits of uint80.
func (u *uint80) setLo(lo uint64) {
	// converts to bigendian
	u[2] = byte(lo >> 56)
	u[3] = byte(lo >> 48)
	u[4] = byte(lo >> 40)
	u[5] = byte(lo >> 32)
	u[6] = byte(lo >> 24)
	u[7] = byte(lo >> 16)
	u[8] = byte(lo >> 8)
	u[9] = byte(lo)
}

// add will increment uint80 by 'i', returning false on overflow.
func (u *uint80) add(i uint64) bool {
	// Load values
	lo := u.lo()
	hi := u.hi()

	// Temp copies
	tlo := lo
	thi := hi

	// Update the uint80 parts
	if tlo += i; tlo < lo {
		thi++
	}

	// Check for overflow
	if thi < hi {
		return false
	}

	// Update self
	u.setHi(thi)
	u.setLo(tlo)

	return true
}
