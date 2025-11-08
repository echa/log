// Copyright (c) 2018-2022 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"io"
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
	Trace(...any)
	Tracef(string, ...any)
	Debug(...any)
	Debugf(string, ...any)
	Info(...any)
	Infof(string, ...any)
	Warn(...any)
	Warnf(string, ...any)
	Error(...any)
	Errorf(string, ...any)
	Fatal(...any)
	Fatalf(string, ...any)
	Panic(...any)
	Panicf(string, ...any)
	Level() Level
	IsColor() bool
	SetLevel(Level) Logger
	SetLevelString(string) Logger
	Logger() *stdlog.Logger
	Clone(string) Logger
	WithTag(string) Logger
	WithRegistry(*Registry) Logger
	WithSampler(*Sampler) Logger
	WithColor(bool) Logger
	WithFlags(int) Logger
	Attach(io.Writer)
	Detach(io.Writer)
}

// package level forwarders to the real logger implementation
func Trace(v ...any)                 { Log.Trace(v...) }
func Tracef(s string, v ...any)      { Log.Tracef(s, v...) }
func Error(v ...any)                 { Log.Error(v...) }
func Errorf(s string, v ...any)      { Log.Errorf(s, v...) }
func Warn(v ...any)                  { Log.Warn(v...) }
func Warnf(s string, v ...any)       { Log.Warnf(s, v...) }
func Info(v ...any)                  { Log.Info(v...) }
func Infof(s string, v ...any)       { Log.Infof(s, v...) }
func Debug(v ...any)                 { Log.Debug(v...) }
func Debugf(s string, v ...any)      { Log.Debugf(s, v...) }
func Fatal(v ...any)                 { Log.Fatal(v...) }
func Fatalf(s string, v ...any)      { Log.Fatalf(s, v...) }
func Panic(v ...any)                 { Log.Panic(v...) }
func Panicf(s string, v ...any)      { Log.Panicf(s, v...) }
func SetLevel(l Level) Logger        { Log.SetLevel(l); return Log }
func SetLevelString(l string) Logger { return SetLevel(ParseLevel(l)) }
