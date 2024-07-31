// Copyright (c) 2018-2022 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var defaultFlags int = stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds | stdlog.LUTC

var levelStrs = [...]string{"TRCE ", "DEBG ", "INFO ", "WARN ", "ERRO ", "CRIT ", "OFF  "}

var levelColors = [...]*color.Color{
	color.New(ColorTrace),
	color.New(ColorDebug),
	color.New(ColorInfo),
	color.New(ColorWarn),
	color.New(ColorError),
	color.New(ColorFatal),
}

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	case "off":
		return LevelOff
	default:
		return LevelInvalid
	}
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelOff:
		return "off"
	default:
		return ""
	}
}

func (l Level) Prefix() string {
	if l >= LevelOff {
		return "off"
	}
	return levelStrs[l]
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(data []byte) error {
	v := ParseLevel(string(data))
	if v == LevelInvalid {
		return fmt.Errorf("invalid log level '%s'", string(data))
	}
	*l = v
	return nil
}

type Config struct {
	Level            Level         `json:"level"`
	Flags            int           `json:"flags"`
	Backend          string        `json:"backend"`
	Addr             string        `json:"addr"`
	Facility         string        `json:"facility"`
	Ident            string        `json:"ident"`
	Filename         string        `json:"filename"`
	FileMode         os.FileMode   `json:"filemode"`
	ProgressInterval time.Duration `json:"progress"`
	NoColor          bool          `json:"nocolor"`
}

func NewConfig() *Config {
	c := &Config{
		Level:            LevelInfo,
		Flags:            defaultFlags,
		Backend:          "stdout", // stdout, stderr, syslog, file
		Addr:             "",
		Facility:         "local0",
		Ident:            "logfile",
		Filename:         "logfile.log",
		FileMode:         0600,
		ProgressInterval: 10 * time.Second,
	}
	c.ParseEnv()
	return c
}

func ParseFlags(flags string) int {
	var cflags int
	if flags == "" {
		return defaultFlags
	}
	for _, f := range strings.Split(flags, ",") {
		switch f {
		case "longfile":
			cflags |= stdlog.Llongfile
		case "shortfile":
			cflags |= stdlog.Lshortfile
		case "date":
			cflags |= stdlog.Ldate
		case "time":
			cflags |= stdlog.Ltime
		case "micro":
			cflags |= stdlog.Lmicroseconds
		case "utc":
			cflags |= stdlog.LUTC
		}
	}
	return cflags
}

func (cfg *Config) ParseEnv() {
	cfg.Flags = ParseFlags(os.Getenv("LOGFLAGS"))
	cfg.Level = ParseLevel(os.Getenv("LOGLEVEL"))
	if cfg.Level == LevelInvalid {
		cfg.Level = LevelInfo
	}
}

func (cfg *Config) Check() error { return nil }
