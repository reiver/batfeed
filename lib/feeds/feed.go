package libfeeds

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"os"
	"slices"

	"github.com/reiver/go-erorr"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/socialfed-api/lib/order"
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

func (receiver Feed) allFileNames() ([]string, error) {
	var pattern string = libpath.Join(receiver.root, receiver.name, "*.url")

	filenames, err := filepath.Glob(pattern)
	if nil != err {
		var nada []string
		return nada, erorr.Errorf("libfeed: problem getting glob %q: %w", pattern, err)
	}

	return filenames, nil
}

func (receiver Feed) FileNames(order string, limit uint64) ([]string, error) {

	filenames, err := receiver.allFileNames()
	if nil != err {
		var nada []string
		return nada, err
	}

	switch order {
	case liborder.OrderAscending:
		// nothing here
	case liborder.OrderDescending:
		slices.Reverse(filenames)
	default:
		panic(fmt.Sprintf("libfeed: unsupported order: %q", order))
	}

	if int(limit) < len(filenames) {
		filenames = filenames[:int(limit)]
	}

	return filenames, nil
}


// Len returns the number of items in ths feed.
func (receiver Feed) Len() (uint64, error) {

	filenames, err := receiver.allFileNames()
	if nil != err {
		var nada uint64
		return nada, err
	}

	return uint64(len(filenames)), nil
}

func (receiver Feed) Path() string {
	return feedPath(receiver.root, receiver.name)
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
