// Copyright (c) 2018-2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package log

import (
	stdlog "log"
	"os"
)

// no syslog on windows, write to stdout
func NewSyslog(config *Config) *Backend {
	return &Backend{config.Level, stdlog.New(os.Stdout, "", config.Flags), "", nil, false}
}
