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
	if b, ok := Log.(*Backend); ok {
		logger = b.NewLogger(tag)
	} else {
		logger = New(nil).WithTag(tag)
	}
	r.reg[tag] = logger
	return logger
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

	prefix, _, hasWildcard := strings.Cut(tag, wildcard)
	switch {
	case !hasWildcard:
		if log, ok := r.reg[tag]; ok {
			log.SetLevel(lvl)
		}
	case hasWildcard && prefix != "":
		for k, v := range r.reg {
			if strings.HasPrefix(k, prefix) {
				v.SetLevel(lvl)
			}
		}
	}
}
