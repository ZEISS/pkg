package writers

import (
	"io"
	"sync"
)

// PausableWriter is a writer that can be paused and resumed.
type PausableWriter struct {
	writer  io.Writer
	discard bool
	sync.Mutex
}

// NewPausableWriter creates a new PausableWriter.
func NewPausableWriter(w io.Writer) *PausableWriter {
	return &PausableWriter{
		discard: false,
		writer:  w,
	}
}

// Write is the implementation of the io.Writer interface.
func (pw *PausableWriter) Write(p []byte) (n int, err error) {
	if pw.discard {
		return io.Discard.Write(p)
	}

	return pw.writer.Write(p)
}

// Pause is used to pause the writer.
func (pw *PausableWriter) Pause() {
	pw.Lock()
	defer pw.Unlock()
	pw.discard = true
}

// Resume is used to resume the writer.
func (pw *PausableWriter) Resume() {
	pw.Lock()
	defer pw.Unlock()
	pw.discard = false
}
