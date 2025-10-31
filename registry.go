// Copyright (c) 2018-2025 KIDTSUNAMI
// Author: abdul@blockwatch.cc

package log

import (
	"strings"
	"sync"
)

const wildcard = "*"

var r = newRegistry()

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
		return Log.Clone().WithTag(tag)
	}
}

func NewLogger(tag string) Logger {
	return r.New(tag)
}

func RemoveLogger(tag string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reg, tag)
}

func GetLogger(tag string) (Logger, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	logger, ok := r.reg[tag]
	return logger, ok
}

func SetLevels(tag string, lvl Level) {
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
