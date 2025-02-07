package verboten

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/reiver/go-did"
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-json"

	"github.com/reiver/batfeed/srv/http"
	"github.com/reiver/batfeed/srv/log"
)

const path string = "/.well-known/did.json"

func init() {
	log := logsrv.Prefix("www("+path+").init").Begin()
	defer log.End()

	var handler http.Handler = http.HandlerFunc(serveHTTP)

	err := httpsrv.Mux.HandlePath(handler, path)
	if nil != err {
		e := erorr.Errorf("problem registering http-handler with path-mux for path %q: %w", path, err)
		log.Error(e)
		panic(e)
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	log := logsrv.Prefix("www("+path+").serveHTTP").Begin()
	defer log.End()

	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}

	if nil == request {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
		log.Error("nil http-request")
		return
	}

	var method string = request.Method

	if http.MethodGet != method {
		errhttp.ErrHTTPMethodNotAllowed.ServeHTTP(responsewriter, request)
		log.Errorf("bad HTTP method: %q", method)
		return
	}

	var tcpaddr string
	{
		tcpaddr = request.Host
		if "" == tcpaddr {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			log.Errorf("empty tcpaddr (%q)", tcpaddr)
			return
		}
	}

	var host string
	{
		index := strings.LastIndexByte(tcpaddr, ':')
		if 0 <= index {
			host = tcpaddr[:index]
		}
	}

	var diduri string
	{
		const method string = "web"
		var identifier string = host

		thedid, err := did.ConstructDID(method, identifier)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			log.Errorf("problem constructing did-uri with method=%q and identifier=%q: %s", method, identifier, err)
			return
		}

		diduri = thedid.String()
	}

	var serviceEndpoint string = fmt.Sprintf("https://%s", tcpaddr)

	var bytes []byte
	{
		type service struct {
			ID              string `json:"id"`
			ServiceEndpoint string `json:"serviceEndpoint"`
			Type            string `json:"type"`
		}

		response := struct {
			Context []string  `json:"@context"`
			ID        string  `json:"id"`
			Service []service `json:"service"`
		}{
			Context: []string{"https://www.w3.org/ns/did/v1"},
			ID: diduri,
			Service: []service{
				service{
					ID: "#bsky_fg",
					ServiceEndpoint: serviceEndpoint,
					Type: "BskyFeedGenerator",
				},
			},
		}

		var err error
		bytes, err = json.Marshal(response)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			log.Errorf("problem marshaling JSON: %s", err)
			return
		}
	}

	{
		responsewriter.Header().Add("Content-Type", "application/json")

		_, err := responsewriter.Write(bytes)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			log.Errorf("problem sending bytes to client: %s", err)
			return
		}
	}
}
