package verboten

import (
	"net/http"
	liburl "net/url"

	"github.com/reiver/go-act"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-jsonld"
	"github.com/reiver/go-opt"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/batfeed/lib/feeds"
	"github.com/reiver/batfeed/srv/feed"
	"github.com/reiver/batfeed/srv/http"
	"github.com/reiver/batfeed/srv/log"
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

	var host string = httprequest.Host

	var httprequesturl *liburl.URL = httprequest.URL
	if nil == httprequesturl {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Error("nil http-request-url")
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

	var bytes []byte
	{
		var actOutbox act.Outbox
		{
			var first = liburl.URL{
				Scheme: "https",
				Host:   host,
				Path:   libpath.Join(httprequesturl.Path, "page"),
				RawQuery: "order=descending",
			}

			var last = liburl.URL{
				Scheme: "https",
				Host:   host,
				Path:   libpath.Join(httprequesturl.Path, "page"),
				RawQuery: "order=ascending",
			}

			var outbox = liburl.URL{
				Scheme: "https",
				Host:   host,
				Path:   httprequesturl.Path,
			}

			var firstString string = first.String()

			var totalItems uint64
			{
				var err error
				totalItems, err = feed.Len()
				if nil != err {
					errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
					log.Errorf("problem getting length of feed %q: %s", feedname, err)
					return
				}
			}
			log.Debugf("total-items: %d", totalItems)

			actOutbox.Current     = opt.Something(firstString)
			actOutbox.First       = opt.Something(firstString)
			actOutbox.ID          = opt.Something(outbox.String())
			actOutbox.Last        = opt.Something(last.String())
			actOutbox.TotalItems  = opt.Something(totalItems)
		}

		var err error
		bytes, err = jsonld.Marshal(actOutbox)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Errorf("problem marshaling json-ld: %s", err)
			return
		}
	}

	act.ServeActivity(responsewriter, request.HTTPRequest(), bytes)
}
