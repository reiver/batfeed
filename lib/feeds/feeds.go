package libfeeds

import (
	"io/fs"
	"os"

	"github.com/reiver/go-erorr"
	libpath "github.com/reiver/go-path"
)

var (
	ErrNotFound error = erorr.Error("not found")
)

const (
	errNilFile       = erorr.Error("nil file")
	errNilFileInfo   = erorr.Error("nil file-info")
	errNilFileSystem = erorr.Error("nil file-system")
)

// Feeds represents an files & directories style database that represents feeds and each of their contents.
type Feeds struct {
	path string
	filesystem fs.FS
}

// OpenFeeds often the directory provided by `path` as a "feeds" directory.
//
// A "feeds" directory is expected to have certain contents in it structured in a certain way.
//
// If the directory given by `path` does not exist, OpenFeeds returns an [ErrNotFound] error.
func OpenFeeds(path string) (Feeds, error) {
	filesystem := os.DirFS(path)
	if nil == filesystem {
		var nada Feeds
		return nada, errNilFileSystem
	}

	return Feeds{
		path:path,
		filesystem:filesystem,
	}, nil
}

// CreateFeed creates a new feed whose name is given by `name`.
func (receiver Feeds) CreateFeed(name string) error {
	var path string = libpath.Join(receiver.path, name)

	err := os.Mkdir(path, 0755)
	if nil != err {
		return erorr.Errorf("problem creating feed %q in %q: %w", name, receiver.path, err)
	}

	return nil
}

// OpenFeed opens a feed whose name is given by `name`, and returns it as a [Feed].
func (receiver Feeds) OpenFeed(name string) (Feed, error) {
	if nil == receiver.filesystem {
		var nada Feed
		return nada, errNilFileSystem
	}

	{
		var path string = libpath.Join(receiver.path, name)

		file, err := os.Open(path)
		if nil != err {
			switch {
			case erorr.Is(err, fs.ErrNotExist) == true:
				var nada Feed
				return nada, ErrNotFound
			default:
				var nada Feed
				return nada, erorr.Errorf("problem opening feed %q at %q: %w", name, receiver.path, err)
			}
		}
		if nil == file {
			var nada Feed
			return nada, errNilFile
		}

		fileinfo, err := file.Stat()
		if nil != err {
			var nada Feed
			return nada, erorr.Errorf("problem getting file-info for %q at %q: %w", name, receiver.path, err)
		}
		if nil == fileinfo {
			var nada Feed
			return nada, errNilFileInfo
		}

		if !fileinfo.IsDir() {
			var nada Feed
			return nada, erorr.Errorf("%q at %q is not a directory", name, receiver.path)
		}
	}

	subfs, err := fs.Sub(receiver.filesystem, name)
	if nil != err {
		var nada Feed
		return nada, erorr.Errorf("problem getting sub-directory %q of %q: %w", name, receiver.path, err)
	}

	return Feed{
		root: receiver.path,
		name: name,
		filesystem: subfs,
	}, nil
}

// FS returns the underlying [fs.FS] of the Feeds.
//
// NOTE, THIS WILL PROBABLY BE REMOVED LATER.
func (receiver Feeds) FS() fs.FS {
	return receiver.filesystem
}
