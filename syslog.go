// Copyright (c) 2018-2022 KIDTSUNAMI
// Author: alex@kidtsunami.com
//go:build !windows
// +build !windows

package log

import (
	stdlog "log"
	"log/syslog"
	"strings"
)

func NewSyslog(c *Config) *Backend {
	if c.Addr != "" {
		parts := strings.Split(c.Addr, "://")
		if len(parts) != 2 {
			stdlog.Fatalln("FATAL: Invalid syslog address. Must be of form protocol://path (e.g. unix:///dev/log)")
		}

		writer, err := syslog.Dial(
			parts[0],
			parts[1],
			syslogFacilityToEnum(c.Facility)|syslog.LOG_INFO,
			c.Ident,
		)
		if err != nil {
			stdlog.Fatalln("FATAL: Cannot open syslog address", c.Addr, ":", err.Error())
		}
		// don't 'print' date time
		return &Backend{c.Level, stdlog.New(writer, "", 0), "", nil, false}
	} else {
		writer, err := syslog.New(
			syslogFacilityToEnum(c.Facility)|syslog.LOG_INFO,
			c.Ident,
		)
		if err != nil {
			stdlog.Fatalln("FATAL: Cannot open syslog:", err.Error())
		}
		// don't 'print' date time
		return &Backend{c.Level, stdlog.New(writer, "", 0), "", nil, false}
	}
}

func syslogFacilityToEnum(f string) (p syslog.Priority) {
	switch strings.ToLower(f) {
	case "kern":
		p = syslog.LOG_KERN
	case "user":
		p = syslog.LOG_USER
	case "mail":
		p = syslog.LOG_MAIL
	case "daemon":
		p = syslog.LOG_DAEMON
	case "auth":
		p = syslog.LOG_AUTH
	case "syslog":
		p = syslog.LOG_SYSLOG
	case "lpr":
		p = syslog.LOG_LPR
	case "news":
		p = syslog.LOG_NEWS
	case "uucp":
		p = syslog.LOG_UUCP
	case "cron":
		p = syslog.LOG_CRON
	case "authpriv":
		p = syslog.LOG_AUTHPRIV
	case "ftp":
		p = syslog.LOG_FTP
	case "local0":
		p = syslog.LOG_LOCAL0
	case "local1":
		p = syslog.LOG_LOCAL1
	case "local2":
		p = syslog.LOG_LOCAL2
	case "local3":
		p = syslog.LOG_LOCAL3
	case "local4":
		p = syslog.LOG_LOCAL4
	case "local5":
		p = syslog.LOG_LOCAL5
	case "local6":
		p = syslog.LOG_LOCAL6
	case "local7":
		p = syslog.LOG_LOCAL7
	default:
		stdlog.Fatalln("FATAL: Invalid syslog facility", f)
	}
	return
}
