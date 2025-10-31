// Copyright (c) 2018-2025 KIDTSUNAMI
// Author: abdul@blockwatch.cc

package log

import (
	"strings"
	"sync"
)

const wildcard = "*"

var (
	r            = newRegistry()
	NewLogger    = r.New
	GetLogger    = r.Get
	RemoveLogger = r.Remove
	SetLevels    = r.SetLevels
)

type registry struct {
	mu  sync.RWMutex
	reg map[string]Logger
}

func newRegistry() *registry {
	return &registry{
		reg: make(map[string]Logger),
	}
}

func (r *registry) New(tag string) (logger Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if logger, ok := r.reg[tag]; ok {
		return logger
	} else {
		logger = Log.Clone().WithTag(tag)
		r.reg[tag] = logger
		return logger
	}
}

func (r *registry) Remove(tag string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reg, tag)
}

func (r *registry) Get(tag string) (Logger, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	logger, ok := r.reg[tag]
	return logger, ok
}

func (r *registry) SetLevels(tag string, lvl Level) {
	r.mu.Lock()
	defer r.mu.Unlock()

	prefix, suffix, hasWildcard := strings.Cut(tag, wildcard)
	switch {
	case !hasWildcard:
		if log, ok := r.reg[tag]; ok {
			log.SetLevel(lvl)
		}
	case hasWildcard:
		for k, v := range r.reg {
			if strings.HasPrefix(k, prefix) && strings.HasSuffix(k, suffix) {
				v.SetLevel(lvl)
			}
		}
	}
}
