package feedsrv

import (
	"fmt"

	"github.com/reiver/batfeed/cfg"
	"github.com/reiver/batfeed/lib/feeds"
)

var Feeds libfeeds.Feeds

func init() {
	var path string = cfg.FeedsDirectory()

	var err error
	Feeds, err = libfeeds.OpenFeeds(path)
	if nil != err {
		panic(fmt.Sprintf("problem opening feeds %q: %s", path, err))
	}
}
