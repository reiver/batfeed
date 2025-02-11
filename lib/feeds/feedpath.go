package libfeeds

import (
	libpath "github.com/reiver/go-path"
)

func feedPath(root string, name string) string {
	if "" == root {
		return ""
	}
	if "" == name {
		return ""
	}

	return libpath.Canonical(libpath.Join(root, name))
}
