package verboten

import (
	"net/http"
	liburl "net/url"

	"github.com/reiver/go-asns"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-jsonld"
	"github.com/reiver/go-opt"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/socialfed/lib/feeds"
	"github.com/reiver/socialfed/lib/order"
	"github.com/reiver/socialfed/srv/http"
	"github.com/reiver/socialfed/srv/log"
)

func serveGET(responsewriter http.ResponseWriter, request *httpsrv.ParameterizedRequest, feedname string, feed libfeeds.Feed) {
	log := logsrv.Prefix("www("+path+").serveGET").Begin()
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
	log.Debugf("host = %q", host)

	var httprequesturl *liburl.URL = httprequest.URL
	if nil == httprequesturl {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Error("nil http-request-url")
		return
	}

	var bytes []byte
	{
		var actOutbox asns.Outbox
		{
			var first = liburl.URL{
				Scheme: "https",
				Host:   host,
				Path:   libpath.Join(httprequesturl.Path, "page"),
				RawQuery: "order=" + liborder.OrderDescending,
			}

			var last = liburl.URL{
				Scheme: "https",
				Host:   host,
				Path:   libpath.Join(httprequesturl.Path, "page"),
				RawQuery: "order=" + liborder.OrderAscending,
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

	asns.ServeActivity(responsewriter, request.HTTPRequest(), bytes)
}
