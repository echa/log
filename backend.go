// Copyright (c) 2018-2022 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
)

type Backend struct {
	level    Level
	log      *log.Logger
	tag      string
	sampler  *Sampler
	usecolor bool
}

var (
	Log         Logger = New(NewConfig())
	Disabled    Logger = &Backend{level: LevelOff, log: log.Default()}
	isColorTerm bool   = !color.NoColor
)

const calldepth = 4

func Init(config *Config) {
	Log = New(config)
}

func New(config *Config) *Backend {
	defaultProgressInterval = config.ProgressInterval
	withColor := isColorTerm && !config.NoColor
	switch strings.ToLower(config.Backend) {
	case "file":
		if config.Filename != "" {
			if file, err := os.OpenFile(config.Filename,
				os.O_WRONLY|os.O_CREATE|os.O_APPEND, config.FileMode); err == nil {
				return &Backend{config.Level, log.New(file, "", config.Flags), "", nil, false}
			} else {
				log.Fatalln("FATAL: Cannot open logfile", config.Filename, ":", err.Error())
			}
		}
	case "syslog":
		return NewSyslog(config)
	case "stdout":
		return &Backend{config.Level, log.New(os.Stdout, "", config.Flags), "", nil, withColor}
	case "stderr":
		return &Backend{config.Level, log.New(os.Stderr, "", config.Flags), "", nil, withColor}
	default:
		log.Fatalln("FATAL: Invalid log backend", config.Backend)
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
	return x
}

func (x *Backend) WithFlags(f int) Logger {
	x.log.SetFlags(f)
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

func (x Backend) Logger() *log.Logger {
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
	}
	m := append(make([]any, 0, len(v)+2), lvl.Prefix(), x.tag)
	m = append(m, v...)
	if x.usecolor {
		x.log.Output(calldepth, levelColors[lvl].Sprint(m...))
	} else {
		x.log.Output(calldepth, fmt.Sprint(m...))
	}
}

func (x Backend) outputf(lvl Level, f string, v ...any) {
	f = strings.Join([]string{"%s%s", f}, "") // prefix tag and level %s
	m := append(make([]any, 0, len(v)+2), lvl.Prefix(), x.tag)
	m = append(m, v...)
	if x.usecolor {
		x.log.Output(calldepth, levelColors[lvl].Sprintf(f, m...))
	} else {
		x.log.Output(calldepth, fmt.Sprintf(f, m...))
	}
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
