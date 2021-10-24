// Copyright (c) 2021 KIDTSUNAMI
// Author: alex@kidtsunami.com
// inspired by https://github.com/rs/zerolog

package log

import (
	"sync/atomic"
	"time"
)

var (
	// Often samples log 10 events per second.
	SampleMany = &Sampler{N: 10, Period: time.Second}
	// Sometimes samples log 1 event per second.
	SampleSome = &Sampler{N: 1, Period: time.Second}
	// Rarely samples log 1 events per minute.
	SampleFew = &Sampler{N: 1, Period: time.Minute}
)

// Sampler lets a burst of N events pass per Period. If Period is0,
// every Nth event is allowed.
type Sampler struct {
	// N is the maximum number of event per period allowed.
	N uint32
	// Period defines the period.
	Period time.Duration

	counter uint32
	resetAt int64
}

func (s *Sampler) Clone() *Sampler {
	if s == nil {
		return s
	}
	return &Sampler{
		N:      s.N,
		Period: s.Period,
	}
}

func (s *Sampler) Sample() bool {
	if s.Period == 0 {
		return s.sampleSimple()
	}
	return s.samplePeriod()
}

func (s *Sampler) sampleSimple() bool {
	n := s.N
	if n == 1 {
		return true
	}
	c := atomic.AddUint32(&s.counter, 1)
	return c%n == 1
}

func (s *Sampler) samplePeriod() bool {
	if s.N > 0 && s.Period > 0 {
		if s.inc() <= s.N {
			return true
		}
	}
	return false
}

func (s *Sampler) inc() uint32 {
	now := time.Now().UnixNano()
	resetAt := atomic.LoadInt64(&s.resetAt)
	var c uint32
	if now > resetAt {
		c = 1
		atomic.StoreUint32(&s.counter, c)
		newResetAt := now + s.Period.Nanoseconds()
		reset := atomic.CompareAndSwapInt64(&s.resetAt, resetAt, newResetAt)
		if !reset {
			// Lost the race with another goroutine trying to reset.
			c = atomic.AddUint32(&s.counter, 1)
		}
	} else {
		c = atomic.AddUint32(&s.counter, 1)
	}
	return c
}
