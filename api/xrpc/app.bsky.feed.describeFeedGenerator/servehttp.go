package verboten

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/reiver/go-did"
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-json"

	"github.com/reiver/batfeed/srv/db"
	"github.com/reiver/batfeed/srv/http"
	. "github.com/reiver/batfeed/srv/log"
)

const path string = "/xrpc/app.bsky.feed.describeFeedGenerator"

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

	var host string
	{
		host = request.Host
		if "" == host {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			Logf("[serve-http][path=%q] empty host (%q)", path, host)
			return
		}
	}

	var domain string
	{
		index := strings.LastIndexByte(host, ':')
		if 0 <= index {
			domain = host[:index]
		}
	}

	var thedid string
	{
		const method string = "web"
		var identifier string = domain

		didObj, err := did.ConstructDID(method, identifier)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			Logf("[serve-http][path=%q] problem constructing did-uri with method=%q and identifier=%q: %s", path, method, identifier, err)
			return
		}

		thedid = didObj.String()
	}

	type feedType struct {
		URI string `json:"uri"`
	}

	var feeds []feedType
	{
		names, err := dbsrv.Feeds(domain)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request)
			Logf("[serve-http][path=%q] problem getting feeds for domain=%q: %s", path, domain, err)
			return
		}

		for _, name := range names {
//@TODO: Maybe construct this URI another way.
			var uri string = fmt.Sprintf("at://%s/app.bsky.feed.generator/%s", thedid, name)

			var feed feedType = feedType{
				URI: uri,
			}

			feeds = append(feeds, feed)
		}
	}

	var bytes []byte
	{

		response := struct {
			DID     string   `json:"did"`
			Feeds []feedType `json:feeds`
		}{
			DID: thedid,
			Feeds: feeds,
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
