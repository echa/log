// Copyright (c) 2018-2024 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
)

type Backend struct {
	level    Level
	log      *stdlog.Logger
	tag      string
	sampler  *Sampler
	usecolor bool
}

var (
	Log         Logger = New(NewConfig())
	Disabled    Logger = &Backend{level: LevelOff, log: stdlog.Default()}
	isColorTerm bool   = !color.NoColor
)

const calldepth = 4

func Init(c *Config) {
	Log = New(c)
}

func New(c *Config) *Backend {
	if c == nil {
		c = NewConfig()
	}
	defaultProgressInterval = c.ProgressInterval
	withColor := isColorTerm || os.Getenv("LOGCOLOR") != ""
	color.NoColor = !withColor
	switch strings.ToLower(c.Backend) {
	case "file":
		if c.Filename != "" {
			if file, err := os.OpenFile(c.Filename,
				os.O_WRONLY|os.O_CREATE|os.O_APPEND, c.FileMode); err == nil {
				backend := &Backend{c.Level, stdlog.New(file, "", c.Flags), "", nil, false}
				runtime.SetFinalizer(backend, func(v any) {
					b := v.(*Backend)
					b.log.Writer().(*os.File).Close()
				})
				return backend
			} else {
				stdlog.Fatalln("FATAL: Cannot open logfile", c.Filename, ":", err.Error())
			}
		}
	case "syslog":
		return NewSyslog(c)
	case "stdout":
		return &Backend{c.Level, stdlog.New(os.Stdout, "", c.Flags), "", nil, withColor}
	case "stderr":
		return &Backend{c.Level, stdlog.New(os.Stderr, "", c.Flags), "", nil, withColor}
	default:
		stdlog.Fatalln("FATAL: Invalid log backend", c.Backend)
	}
	return nil
}

func (x Backend) NewLogger(subsystem string) Logger {
	return &Backend{
		level:    x.level,
		log:      x.log,
		tag:      strings.TrimSpace(subsystem) + " ",
		usecolor: x.usecolor,
	}
}

func (x Backend) Clone() Logger {
	return &Backend{
		level:    x.level,
		log:      x.log,
		tag:      x.tag,
		sampler:  x.sampler.Clone(),
		usecolor: x.usecolor,
	}
}

func (x *Backend) WithTag(tag string) Logger {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		x.tag += tag + " "
	}
	return x
}

func (x *Backend) WithSampler(s *Sampler) Logger {
	x.sampler = s
	return x
}

func (x *Backend) WithColor(b bool) Logger {
	x.usecolor = b
	color.NoColor = !b
	return x
}

func (x *Backend) IsColor() bool {
	return x.usecolor
}

func (x *Backend) WithFlags(f int) Logger {
	x.log.SetFlags(f)
	return x
}

func (x *Backend) WithLogger(l *stdlog.Logger) Logger {
	x.log = l
	return x
}

func (x Backend) NewWriter(l Level) io.Writer {
	if x.level > l {
		return ioutil.Discard
	}
	writer := &Backend{
		level:    l,
		log:      x.log,
		tag:      x.tag,
		usecolor: x.usecolor,
	}
	return writer
}

func (x Backend) Write(p []byte) (n int, err error) {
	if l := len(p); l == 0 {
		return 0, nil
	} else if p[l-1] == '\n' {
		x.output(x.level, string(p[:l-1]))
		return l - 1, nil
	} else {
		x.output(x.level, string(p))
		return l, nil
	}
}

func (x Backend) Logger() *stdlog.Logger {
	return x.log
}

func (x Backend) Level() Level {
	return x.level
}

func (x *Backend) SetLevel(l Level) Logger {
	if l != LevelInvalid {
		x.level = l
	}
	return x
}

func (x *Backend) SetLevelString(s string) Logger {
	return x.SetLevel(ParseLevel(s))
}

func (x Backend) Noop(...any) {}

func (x Backend) Error(v ...any) {
	if !x.shouldLog(LevelError) {
		return
	}
	x.output(LevelError, v...)
}

func (x Backend) Errorf(f string, v ...any) {
	if !x.shouldLog(LevelError) {
		return
	}
	x.outputf(LevelError, f, v...)
}

func (x Backend) Warn(v ...any) {
	if !x.shouldLog(LevelWarn) {
		return
	}
	x.output(LevelWarn, v...)
}

func (x Backend) Warnf(f string, v ...any) {
	if !x.shouldLog(LevelWarn) {
		return
	}
	x.outputf(LevelWarn, f, v...)
}

func (x Backend) Info(v ...any) {
	if !x.shouldLog(LevelInfo) {
		return
	}
	x.output(LevelInfo, v...)
}

func (x Backend) Infof(f string, v ...any) {
	if !x.shouldLog(LevelInfo) {
		return
	}
	x.outputf(LevelInfo, f, v...)
}

func (x Backend) Debug(v ...any) {
	if !x.shouldLog(LevelDebug) {
		return
	}
	x.output(LevelDebug, v...)
}

func (x Backend) Debugf(f string, v ...any) {
	if !x.shouldLog(LevelDebug) {
		return
	}
	x.outputf(LevelDebug, f, v...)
}

func (x Backend) Fatal(v ...any) {
	x.output(LevelFatal, v...)
	x.stackTrace(LevelFatal, 3)
	x.output(LevelFatal, "Exiting process")
	os.Exit(1)
}

func (x Backend) Fatalf(f string, v ...any) {
	x.outputf(LevelFatal, f, v...)
	x.stackTrace(LevelFatal, 3)
	x.output(LevelFatal, "Exiting process")
	os.Exit(1)
}

func (x Backend) Trace(v ...any) {
	if !x.shouldLog(LevelTrace) {
		return
	}
	x.output(LevelTrace, v...)
}

func (x Backend) Tracef(f string, v ...any) {
	if !x.shouldLog(LevelTrace) {
		return
	}
	x.outputf(LevelTrace, f, v...)
}

func (x Backend) output(lvl Level, v ...any) {
	if len(v) == 1 {
		if fn, ok := v[0].(func()); ok {
			fn()
			return
		}
		if fn, ok := v[0].(func() string); ok {
			v[0] = fn()
		}
	}
	m := append(make([]any, 0, len(v)+2), lvl.Prefix(), x.tag)
	m = append(m, v...)
	print := fmt.Sprint
	if x.usecolor {
		print = levelColors[lvl].Sprint
	}
	x.log.Output(calldepth, print(m...))
}

func (x Backend) outputf(lvl Level, f string, v ...any) {
	f = strings.Join([]string{"%s%s", f}, "") // prefix tag and level %s
	m := append(make([]any, 0, len(v)+2), lvl.Prefix(), x.tag)
	m = append(m, v...)
	print := fmt.Sprintf
	if x.usecolor {
		print = levelColors[lvl].Sprintf
	}
	x.log.Output(calldepth, print(f, m...))
}

func (x Backend) shouldLog(lvl Level) bool {
	if x.level > lvl {
		return false
	}
	if x.sampler != nil {
		return x.sampler.Sample()
	}
	return true
}

func (x Backend) stackTrace(lvl Level, skip int) {
	trace := debug.Stack()
	skip = skip*2 + 1
	for _, v := range bytes.Split(trace, []byte("\n")) {
		if len(v) == 0 {
			continue
		}
		if skip > 0 {
			skip--
			continue
		}
		x.output(lvl, string(v))
	}
}
