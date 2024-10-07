// Copyright (c) 2018-2022 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	stdlog "log"

	"github.com/fatih/color"
)

type LogFn func(...any)
type LogfFn func(string, ...any)
type Level byte

var Noop = func(string, ...any) {}

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
	LevelInvalid
)

var (
	ColorTrace = color.FgHiBlue
	ColorDebug = color.FgCyan
	ColorInfo  = color.FgGreen
	ColorWarn  = color.FgYellow
	ColorError = color.FgRed
	ColorFatal = color.FgMagenta
)

type Logger interface {
	Noop(...any)
	Trace(v ...any)
	Tracef(f string, v ...any)
	Debug(v ...any)
	Debugf(f string, v ...any)
	Info(v ...any)
	Infof(f string, v ...any)
	Warn(v ...any)
	Warnf(f string, v ...any)
	Error(v ...any)
	Errorf(f string, v ...any)
	Fatal(v ...any)
	Fatalf(f string, v ...any)
	Level() Level
	IsColor() bool
	SetLevel(Level) Logger
	SetLevelString(string) Logger
	Logger() *stdlog.Logger
	Clone() Logger
	WithTag(tag string) Logger
	WithSampler(s *Sampler) Logger
	WithColor(b bool) Logger
	WithFlags(f int) Logger
	WithLogger(l *stdlog.Logger) Logger
}

// package level forwarders to the real logger implementation
func Trace(v ...any)            { Log.Trace(v...) }
func Tracef(s string, v ...any) { Log.Tracef(s, v...) }
func Error(v ...any)            { Log.Error(v...) }
func Errorf(s string, v ...any) { Log.Errorf(s, v...) }
func Warn(v ...any)             { Log.Warn(v...) }
func Warnf(s string, v ...any)  { Log.Warnf(s, v...) }
func Info(v ...any)             { Log.Info(v...) }
func Infof(s string, v ...any)  { Log.Infof(s, v...) }
func Debug(v ...any)            { Log.Debug(v...) }
func Debugf(s string, v ...any) { Log.Debugf(s, v...) }
func Fatal(v ...any)            { Log.Fatal(v...) }
func Fatalf(s string, v ...any) { Log.Fatalf(s, v...) }

func SetLevel(l Level) Logger { Log.SetLevel(l); return Log }

func SetLevelString(l string) Logger { return SetLevel(ParseLevel(l)) }

func NewLogger(tag string) Logger {
	if b, ok := Log.(*Backend); ok {
		return b.NewLogger(tag)
	} else {
		return New(nil).WithTag(tag)
	}
}
