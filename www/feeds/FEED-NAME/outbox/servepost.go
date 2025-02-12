package verboten

import (
	gojson "encoding/json"
	"io"
	"net/http"
	liburl "net/url"

	"github.com/reiver/go-asns"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-http201"
	"github.com/reiver/go-http400"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/socialfed-api/lib/feeds"
	"github.com/reiver/socialfed-api/lib/urlfile"
	"github.com/reiver/socialfed-api/srv/http"
	"github.com/reiver/socialfed-api/srv/log"
)

func servePOST(responsewriter http.ResponseWriter, request *httpsrv.ParameterizedRequest, feedname string, feed libfeeds.Feed) {
	log := logsrv.Prefix("www("+path+").servePOST").Begin()
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

	var httprequesturl *liburl.URL = httprequest.URL
	if nil == httprequesturl {
		errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
		log.Error("nil http-request-url")
		return
	}

	var host string = httprequest.Host
	log.Debugf("host = %q", host)

	var contentType string
	{
		var header http.Header = httprequest.Header
		if nil == header {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Debug("nil http-request header")
			return
		}

		contentType = header.Get("Content-Type")
		if "" == contentType {
			http400.BadRequest(responsewriter, request.HTTPRequest())
			log.Debug("empty Content-Type")
			return
		}
	}

	var objecturl string
	switch contentType {
	case `application/x-www-form-urlencoded`:
		err := httprequest.ParseForm()
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Debugf("problem parsing form: %s", err)
			return
		}

		var form liburl.Values = httprequest.PostForm
		if nil == form {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Debug("nil form")
			return
		}

		if "" == objecturl {
			objecturl = form.Get("object")
		}
		if "" == objecturl {
			objecturl = form.Get("ref")
		}
		if "" == objecturl {
			objecturl = form.Get("uri")
		}
		if "" == objecturl {
			objecturl = form.Get("url")
		}
		if "" == objecturl {
			http400.BadRequest(responsewriter, request.HTTPRequest())
			log.Debugf("bad object url: %s", objecturl)
			return
		}
	case `application/activity+json`, `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`:
		var bytes []byte
		{
			var err error
			bytes, err = io.ReadAll(httprequest.Body)
			if nil != err {
				errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
				log.Error("problem reading-all http-request-body")
				return
			}
			httprequest.Body.Close()
		}

		var data map[string]string = map[string]string{}

		err := gojson.Unmarshal(bytes, &data)
		if nil != err {
			http400.BadRequest(responsewriter, request.HTTPRequest())
			log.Debugf("problem unmarshaling JSON: %s", err)
			return
		}

		var typefield string
		{
			var found bool
			typefield, found = data["type"]
			if !found {
				http400.BadRequest(responsewriter, request.HTTPRequest())
				log.Debug("missing type")
				return
			}
			if asns.ActivityTypeAnnounce != typefield {
				http400.BadRequest(responsewriter, request.HTTPRequest())
				log.Debug("missing type")
				return
			}
		}

		{
			var found bool
			objecturl, found = data["object"]
			if !found {
				http400.BadRequest(responsewriter, request.HTTPRequest())
				log.Debug("missing object-url")
				return
			}
		}
	default:
			http400.BadRequest(responsewriter, request.HTTPRequest())
			log.Debugf("unsupported Content-Type: %q", contentType)
			return
	}
	log.Debugf("object-url = %q", objecturl)
	{
		urloc, err := liburl.Parse(objecturl)
		if nil != err {
			http400.BadRequest(responsewriter, request.HTTPRequest())
			log.Debugf("problem parsing url %q: %s", objecturl, err)
			return
		}
		if nil == urloc {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Error("nil parsed-url")
			return
		}

		switch urloc.Scheme {
		case "http","https":
			// nothing here
		default:
			http400.BadRequest(responsewriter, request.HTTPRequest())
			log.Debugf("problem parsing url %q: %s", objecturl, err)
			return
		}
	}

	var internalFileName string
	{
		var err error
		internalFileName, err = feed.Post(objecturl)
		if nil != err {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Errorf("problem posting object-URL %q to feed %q: %s", objecturl, feedname, err)
			return
		}

	}
	log.Debugf("internal-file-name: %q", internalFileName)

	var location string
	{
		var dir string = libpath.Parent(httprequesturl.Path)
		var locpath string = libpath.Join(dir, "objects", liburlfile.Slug(internalFileName))

		if "" == locpath {
			errhttp.ErrHTTPInternalServerError.ServeHTTP(responsewriter, request.HTTPRequest())
			log.Errorf("empty location path: %s", locpath)
			return
		}

		var loc liburl.URL
		loc.Scheme = "https"
		loc.Host = host
		loc.Path = locpath

		location = loc.String()
	}
	log.Debugf("location: %q", location)

	{
		err := http201.ServeLocation(responsewriter, location)
		if nil != err {
			log.Errorf("problem serving http-201 (created) response with location %q: %s", location, err)
			return
		}
	}
}
