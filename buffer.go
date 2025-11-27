// Copyright (c) 2025 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"bytes"
	"iter"
	"sync"
)

// RingBuffer stores capacity messages (always a power of 2). Useful
// with MultiWriter or as general write backend for logging and
// allows retrieval of the last N historic messages.
//
// buf := NewRingBuffer(128)
//
// // replace stdlog backend
// log.Log.Logger().SetOutput(buf)
//
// // attach to current backend
// log.Log.Attach(buf)
type RingBuffer struct {
	mu   sync.RWMutex
	buf  [][]byte
	head uint64
	mask uint64
}

func NewRingBuffer(capacity int) *RingBuffer {
	size := uint64(1)
	for size < uint64(capacity) {
		size <<= 1
	}
	return &RingBuffer{
		buf:  make([][]byte, size),
		mask: size - 1,
	}
}

func (rb *RingBuffer) Write(buf []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	idx := rb.head & rb.mask
	rb.buf[idx] = bytes.Clone(buf)
	rb.head++
	return len(buf), nil
}

func (rb *RingBuffer) Last(n int) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		rb.mu.RLock()
		defer rb.mu.RUnlock()
		n = min(n, len(rb.buf), int(rb.head))
		if n <= 0 {
			return
		}
		idx := (rb.head - uint64(n)) & rb.mask
		for range n {
			if !yield(rb.buf[idx]) {
				return
			}
			idx = (idx + 1) & rb.mask
		}
	}
}
