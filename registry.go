// Copyright (c) 2018-2025 KIDTSUNAMI
// Author: abdul@blockwatch.cc

package log

import (
	"strings"
	"sync"
)

const wildcard = "*"

var (
	DefaultRegistry = NewRegistry()
)

type Registry struct {
	mu  sync.RWMutex
	reg map[string]Logger
}

func NewRegistry() *Registry {
	return &Registry{
		reg: make(map[string]Logger),
	}
}

func (r *Registry) Add(tag string, l Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reg[tag] = l
}

func (r *Registry) Remove(tag string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reg, tag)
}

func (r *Registry) Get(tag string) (Logger, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	logger, ok := r.reg[tag]
	return logger, ok
}

func (r *Registry) GetLevels() map[string]Level {
	r.mu.Lock()
	defer r.mu.Unlock()
	m := make(map[string]Level, len(r.reg))
	for n, l := range r.reg {
		m[n] = l.Level()
	}
	return m
}

func (r *Registry) SetLevels(tag string, lvl Level) {
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
