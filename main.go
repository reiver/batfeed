package main

import (
	"github.com/reiver/batfeed/srv/log"

	// Do this so we get errors early.
	_ "github.com/reiver/batfeed/srv/db"
)

func main() {
	log := logsrv.Prefix("main")
	log.Begin()
	defer log.End()

	log.Inform("BatFeed ⚡")

	log.Inform("Here we go…")
	webserve()
}
