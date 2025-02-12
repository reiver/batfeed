package main

import (
	"github.com/reiver/socialfed/srv/log"

	// Do this so we get errors early.
	_ "github.com/reiver/socialfed/srv/db"
)

func main() {
	log := logsrv.Prefix("main").Begin()
	defer log.End()

	log.Inform("BatFeed ⚡")
	shout()

	log.Inform("Let me show you something…")
	reveal()

	log.Inform("Here we go…")
	webserve()
}
