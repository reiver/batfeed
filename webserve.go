package main

import (
	"net/http"

	"github.com/reiver/batfeed/cfg"
	"github.com/reiver/batfeed/srv/http"
	"github.com/reiver/batfeed/srv/log"
	_ "github.com/reiver/batfeed/www" // this import activates all the HTTP handlers.
)

func webserve() {
	log := logsrv.Prefix("webserve")
	log.Begin()
	defer log.End()

	var tcpaddr string = cfg.WebServerTCPAddress()
	log.Informf("serving HTTP on TCP address: %q", tcpaddr)

	err := http.ListenAndServe(tcpaddr, &httpsrv.Mux)
	if nil != err {
		log.Errorf("ERROR: problem with serving HTTP on TCP address %q: %s", tcpaddr, err)
		panic(err)
	}
}
