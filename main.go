package main

import (
	"fmt"
	"net/http"

	"github.com/reiver/batfeed/srv/http"
	"github.com/reiver/batfeed/srv/log"

	// import these package so their init() fuctions and other initializers run.
	_ "github.com/reiver/batfeed/www"

	// Do this so we get errors early.
	_ "github.com/reiver/batfeed/srv/db"
)

func init() {
	log := logsrv.Prefix("main.init")
	log.Begin()
	defer log.End()

	log.Inform("-<([ hello world ])>-")
	log.Inform()
	log.Inform("batfeed")
	log.Inform()
}

func main() {
	log := logsrv.Prefix("main")
	log.Begin()
	defer log.End()

	var tcpport string = tcpPort()
	log.Informf("tcp-port = %q", tcpport)

	var addr string = fmt.Sprintf(":%s", tcpport)
	log.Informf("tcp-address = %q", addr)

	var handler http.Handler = &httpsrv.Mux

	{
		log.Inform()
		log.Inform("Here we go…")
		err := http.ListenAndServe(addr, handler)
		if nil != err {
			log.Errorf("ERROR: HTTP server had problem listening-and-serving: %s", err)
			return
		}
		log.Inform("beware i live")
	}
}
