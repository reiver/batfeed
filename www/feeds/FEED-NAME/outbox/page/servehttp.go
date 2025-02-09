package verboten

import (
	"fmt"
	"net/http"
	liburl "net/url"
	"strconv"
	"strings"

	"github.com/reiver/go-asns"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-jsonld"
	"github.com/reiver/go-opt"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/batfeed/lib/feeds"
	"github.com/reiver/batfeed/lib/order"
	"github.com/reiver/batfeed/lib/urlfile"
	"github.com/reiver/batfeed/srv/feed"
	"github.com/reiver/batfeed/srv/file"
	"github.com/reiver/batfeed/srv/http"
	"github.com/reiver/batfeed/srv/log"
)

const path string = "/feeds/{feed-name}/outbox/page"

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

	var query liburl.Values = httprequesturl.Query()
	if nil == query {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Error("nil query")
		return
	}

	var order string
	{

		order = query.Get("order")
		if !liborder.IsValidOrder(order) {
			errhttp.ErrHTTPBadRequest.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Errorf("unsupported 'order': %q", order)
			return
		}
	}
	log.Debugf("order: %q", order)

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

	var num uint64
	{
		numStr := query.Get("num")
		switch numStr {
		case "":
			num = 0
		default:
			var err error
			num, err = strconv.ParseUint(numStr, 10, 64)
			if nil != err {
				errhttp.ErrHTTPBadRequest.ServeHTTP(responsewriter, request.HTTPRequest())
				log.Errorf("bad 'num' %q: %s", numStr, err)
				return
			}
		}
	}
	log.Debugf("num = %d", num)

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

	var items []any
	{
		const limit = 20

		filenames, err := feed.FileNames(order, limit)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Errorf("problem opening feed %q: %s", feedname, err)
			return
		}

		var ref liburl.URL = *httprequesturl
		ref.Scheme = "https"
		ref.Host = host
		ref.Fragment = ""
		ref.RawFragment = ""

		var actor liburl.URL = ref
		actor.Path = libpath.RemoveTrailingSeparators(libpath.Parent(libpath.Parent(actor.Path)))
		actor.RawQuery = ""

		var idref liburl.URL = ref
		idref.Path = libpath.Join(libpath.Parent(libpath.Parent(idref.Path)), "objects")
		idref.RawQuery = ""

		for _, filename := range filenames {
			filecontent, err := filesrv.ReadAll(filename)
			if nil != err {
				errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
				log.Errorf("problem reading file %q in feed %q: %s", filename, feedname, err)
				return
			}
			filecontent = strings.TrimSpace(filecontent)

			var id liburl.URL = idref
			id.Path = libpath.Join(id.Path, liburlfile.Slug(filename))

			var published string = liburlfile.Published(filename)

			var announce = asns.Announce{
				Actor:     opt.Something(actor.String()),
				ID:        opt.Something(id.String()),
				Object:    opt.Something(filecontent),
				Published: opt.Something(published),
			}

			items = append(items, announce)
		}
	}

	var bytes []byte
	{
		var ref liburl.URL = *httprequesturl
		ref.Scheme = "https"
		ref.Host = host
		ref.Fragment = ""
		ref.RawFragment = ""

		var id liburl.URL = ref

		var partOf liburl.URL = ref
		partOf.Path = libpath.RemoveTrailingSeparators(libpath.Parent(partOf.Path))
		partOf.RawQuery = ""

		var next liburl.URL = ref
		next.RawQuery = fmt.Sprintf("order=%s&num=%d", order, num+1)

		var asOrderedCollectionPage = asns.OrderedCollectionPage{
			ID:           opt.Something(id.String()),
			OrderedItems: items,
			PartOf:       opt.Something(partOf.String()),
			Next:         opt.Something(next.String()),
		}

		if 0 < num {
			var prev liburl.URL = ref
			if 1 == num {
				prev.RawQuery = fmt.Sprintf("order=%s", order)
			} else {
				prev.RawQuery = fmt.Sprintf("order=%s&num=%d", order, num-1)
			}

			asOrderedCollectionPage.Prev = opt.Something(prev.String())
		}

		var err error
		bytes, err = jsonld.Marshal(asOrderedCollectionPage)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Errorf("problem marshaling json-ld: %s", err)
			return
		}
	}

	asns.ServeActivity(responsewriter, request.HTTPRequest(), bytes)
}
