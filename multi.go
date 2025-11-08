// Copyright (c) 2025 KIDTSUNAMI
// Author: alex@blockwatch.cc

package log

import (
	"io"
	"sync/atomic"
)

// make sure MultiWriter implements io.Writer
var _ io.Writer = (*MultiWriter)(nil)

// MultiWriter is a writer that writes to multiple other writers.
type MultiWriter struct {
	writers atomic.Pointer[[]io.Writer]
}

// New creates a writer that duplicates its writes to all the provided writers,
// similar to the Unix tee(1) command. Writers can be added and removed
// dynamically after creation.
//
// Each write is written to each listed writer, one at a time. Errors returned
// by writers are silently ignored so that a single failed writer does not
// impact others in forwarding log messages.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	mw := &MultiWriter{}
	mw.writers.Store(&writers)
	return mw
}

// Write writes bytes to all writers and silently ignores all errors.
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range *mw.writers.Load() {
		_, _ = w.Write(p)
	}
	return len(p), nil
}

// Add appends a writer to the list of writers this multiwriter writes to.
// Duplicates are igored.
func (mw *MultiWriter) Add(w io.Writer) {
	old := *mw.writers.Load()
	new := make([]io.Writer, 0, len(old)+1)
	for _, ew := range old {
		if ew == w {
			return
		}
		new = append(new, ew)
	}
	new = append(new, w)
	mw.writers.Store(&new)
}

// Remove will remove a previously added writer from the list of writers.
func (mw *MultiWriter) Remove(w io.Writer) {
	var k int
	old := *mw.writers.Load()
	new := make([]io.Writer, len(old))
	for _, ew := range old {
		if ew == w {
			continue
		}
		new[k] = ew
		k++
	}
	new = new[:k]
	mw.writers.Store(&new)
}
