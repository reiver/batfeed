package env

import (
	"os"
)

var FeedsDir string = feedsDir()

// tcpPort returns the directory that the feeds-service should use.
//
// It defaults to directory "./feeds".
//
// But that can be overridden by a value set in the "FEEDS_DIR" environment variable.
func feedsDir() string {
	feedsDir := os.Getenv("FEEDS_DIR")
	if "" == feedsDir {
		feedsDir = "./feeds"
	}

	return feedsDir
}
