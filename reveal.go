package main

import (
	"path/filepath"

	"github.com/reiver/batfeed/cfg"
	"github.com/reiver/batfeed/srv/log"
)

func reveal() {
	log := logsrv.Prefix("reveal").Begin()
	defer log.End()

	var feedsdir string = cfg.FeedsDirectory()
	log.Informf("expected feeds to be in directory: %q", feedsdir)
	{
		abs, err := filepath.Abs(feedsdir)
		if nil == err {
		        log.Informf("absolute-path of directory feeds is expected to be in: %q", abs)
		}
	}

	var tcpaddr string = cfg.WebServerTCPAddress()
	log.Informf("serving HTTP on TCP address: %q", tcpaddr)
}
