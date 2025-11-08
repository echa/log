// Copyright (c) 2018-2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	stdlog "log"
	"os"
)

// no syslog on windows, write to stdout
func NewSyslog(c *Config) *Backend {
	return &Backend{
		level:  c.Level,
		log:    stdlog.New(NewMultiWriter(os.Stdout), "", c.Flags),
		config: c,
	}
}
