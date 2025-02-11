package liburlfile

import (
	"encoding/base64"
	"strings"
	"unsafe"
)

// Slug returns the 'slug' based on the filename part of the `path`.
//
// If `path` contains a directory path, the directory part is removed.
func Slug(path string) string {

	{
		var index int = strings.LastIndex(path, "/")
		if 0 <= index && index < len(path) {
			path = path[1+index:]
		}
		if "" == path {
			return ""
		}
	}

	strdata := unsafe.StringData(path)
	var p []byte = unsafe.Slice(strdata, len(path))

	return base64.RawURLEncoding.EncodeToString(p)
}
