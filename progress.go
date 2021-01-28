// Copyright (c) 2018-2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"strings"
	"sync"
	"time"
)

var defaultProgressInterval = 10 * time.Second

type ProgressLogger struct {
	action      string
	event       string
	interval    time.Duration
	calls       int64
	events      int64
	lastLogTime time.Time
	logger      Logger
	sync.Mutex
}

func NewProgressLogger(logger Logger) *ProgressLogger {
	if logger == nil {
		logger = Log
	}
	return &ProgressLogger{
		action:      "Processed",
		event:       "call",
		interval:    defaultProgressInterval,
		lastLogTime: time.Now(),
		logger:      logger,
	}
}

func (l *ProgressLogger) SetAction(action string) *ProgressLogger {
	l.action = action
	return l
}

func (l *ProgressLogger) SetEvent(event string) *ProgressLogger {
	l.event = event
	return l
}

func (l *ProgressLogger) SetInterval(interval time.Duration) *ProgressLogger {
	l.interval = interval
	return l
}

func pluralize(str string, count int64) string {
	if count == 0 || count > 1 {
		str += "s"
	}
	return str
}

func (p *ProgressLogger) LogN(n int) {
	p.Log(n, time.Time{})
}

func (p *ProgressLogger) Log(n int, ts time.Time, extra ...string) {
	p.Lock()
	defer p.Unlock()
	p.calls++
	p.events += int64(n)
	now := time.Now()
	duration := now.Sub(p.lastLogTime)
	if duration < p.interval || p.events == 0 {
		return
	}

	// Truncate the duration to 10s of milliseconds.
	tDuration := duration.Truncate(10 * time.Millisecond)

	// Log information about the event.
	extraString := ""
	if ex := strings.Join(extra, " "); len(ex) > 0 {
		extraString = " (" + ex + ")"
	}
	tm := ts.UTC().Format("2006-01-02 15:04:05 MST ")
	if ts.IsZero() {
		tm = ""
	}
	p.logger.Infof("%s %d %s in %s%s",
		p.action,
		p.events,
		pluralize(p.event, p.events),
		tDuration,
		tm,
		extraString,
	)

	p.calls = 0
	p.events = 0
	p.lastLogTime = now
}
