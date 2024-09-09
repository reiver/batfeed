package verboten

import (
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-json"

	"github.com/reiver/batfeed/srv/http"
	. "github.com/reiver/batfeed/srv/log"
)

const path string = "/.well-known/did.json"

func init() {
	var handler http.Handler = http.HandlerFunc(serveHTTP)

	err := httpsrv.Mux.HandlePath(handler, path)
	if nil != err {
		e := erorr.Errorf("problem registering http-handler with path-mux for path %q: %w", path, err)
		Log(e)
		panic(e)
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {

	if nil == responsewriter {
		Logf("[serve-http][path=%q] nil http.ResponseWriter", path)
		return
	}

	if nil == request {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
		Logf("[serve-http][path=%q] nil *http.Request", path)
		return
	}

	var method string = request.Method

	if http.MethodGet != method {
		errhttp.ErrHTTPMethodNotAllowed.ServeHTTP(responsewriter, request)
		Logf("[serve-http][path=%q] bad HTTP method: %q", path, method)
		return
	}

	var diduri string
	{
		diduri = "did:web:example.com"
	}

	var serviceEndpoint string
	{
		serviceEndpoint = "http://example.com"
	}

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
			Logf("[serve-http][path=%q] problem marshaling JSON: %s", path, err)
			return
		}
	}

	{
		responsewriter.Header().Add("Content-Type", "application/json")

		_, err := responsewriter.Write(bytes)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			Logf("[serve-http][path=%q] problem sending bytes to client: %s", path, err)
			return
		}
	}
}
