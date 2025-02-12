package verboten

import (
	"net/http"

	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-erorr"

	"github.com/reiver/socialfed/lib/feeds"
	"github.com/reiver/socialfed/srv/feed"
	"github.com/reiver/socialfed/srv/http"
	"github.com/reiver/socialfed/srv/log"
)

const path string = "/feeds/{feed-name}/outbox"

func init() {
	log := logsrv.Prefix("www("+path+").init").Begin()
	defer log.End()

	var handler httpsrv.PatternHandler = httpsrv.PatternHandlerFunc(serveHTTP)

	err := httpsrv.Mux.HandlePattern(handler, path)
	if nil != err {
		e := erorr.Errorf("problem registering http-handler with path-mux for path %q: %w", path, err)
		panic(e)
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *httpsrv.ParameterizedRequest) {
	log := logsrv.Prefix("www("+path+").serveHTTP").Begin()
	defer log.End()

	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Error("nil parameterized-request")
		return
	}

	var httprequest *http.Request = request.HTTPRequest()
	if nil == httprequest {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Error("nil http-request")
		return
	}

	var feedname string
	{
		var found bool
		feedname, found = request.ParameterByIndex(0)
		if !found {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Error("empty feed-name")
			return
		}
	}
	log.Debugf("feedname = %q", feedname)

	var feed libfeeds.Feed
	{
		var err error
		feed, err = feedsrv.Feeds.OpenFeed(feedname)
		if nil != err {
			switch {
			case erorr.Is(err, libfeeds.ErrNotFound):
				errhttp.ErrHTTPNotFound.ServeHTTP(responsewriter, request.HTTPRequest())
				log.Debugf("feed %q not found: %s", feedname, err)
				return
			default:
				errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
				log.Errorf("problem opening feed %q: %s", feedname, err)
				return
			}
			return
		}
	}

	var method string = httprequest.Method
	log.Debugf("method = %q", method)

	switch method {
	case http.MethodGet:
		serveGET(responsewriter, request, feedname, feed)
		return
	case http.MethodPost:
		servePOST(responsewriter, request, feedname, feed)
		return
	default:
		errhttp.ErrHTTPMethodNotAllowed.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Debugf("method not allowed: %q", method)
		return
	}
}
