package libfeeds

import (
	"io/fs"
	"path/filepath"

	"github.com/reiver/go-erorr"
	libpath "github.com/reiver/go-path"
	"os"
)

const (
	errEmptyPath = erorr.Error("empty path")
)

// Feed represent a single feed within a Feeds files & directories  style database.
type Feed struct {
	root string
	name string
	filesystem fs.FS
}

// Len returns the number of items in ths feed.
func (receiver Feed) Len() (uint64, error) {

	filenames, err := filepath.Glob("*.url")
	if nil != err {
		var nada uint64
		return nada, err
	}

	return uint64(len(filenames)), nil
}

// Post adds a new item to a feed.
func (receiver Feed) Post(url string) (string, error) {
	var dirpath string = libpath.Join(receiver.root, receiver.name)
	if "" == dirpath {
		var nada string
		return nada, errEmptyPath
	}

	var filename string = freshURLFileName()
	var path string = libpath.Join(dirpath, filename)

	err := os.WriteFile(path, []byte(url), 0644)
	if nil != err {
		var nada string
		return nada, erorr.Errorf("problem writing file to %q: %w", path, err)
	}

	return filename, nil
}
