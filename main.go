package main

import (
	"github.com/reiver/socialfed-api/srv/log"

	// Do this so we get errors early.
	_ "github.com/reiver/socialfed-api/srv/db"
)

func main() {
	log := logsrv.Prefix("main").Begin()
	defer log.End()

	log.Inform("socialfed ⚡")
	shout()

	log.Inform("Let me show you something…")
	reveal()

	log.Inform("Here we go…")
	webserve()
}
