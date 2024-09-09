package main

import (
	"fmt"
	"net/http"

	"github.com/reiver/batfeed/srv/http"
	. "github.com/reiver/batfeed/srv/log"

	// import these package so their init() fuctions and other initializers run.
	_ "github.com/reiver/batfeed/api"
)

func main() {
	Log("-<([ hello world ])>-")
	Log()
	Log("batfeed")
	Log()

	var tcpport string = tcpPort()
	Logf("tcp-port = %q", tcpport)

	var addr string = fmt.Sprintf(":%s", tcpport)
	Logf("tcp-address = %q", addr)

	var handler http.Handler = &httpsrv.Mux

	{
		Log()
		Log("Here we go…")
		err := http.ListenAndServe(addr, handler)
		if nil != err {
			Logf("ERROR: HTTP server had problem listening-and-serving: %s", err)
			return
		}
		Log("beware i live")
	}
}
