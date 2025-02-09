package liburlfile

import (
	"strings"
)

func Published(path string) string {

	{
		var index int = strings.Index(path, "_")
		if index < 0 {
			return ""
		}
		if len(path) <= index {
			return ""
		}

		path = path[:index]
	}

	{
		var index int = strings.LastIndex(path, "/")
		if 0 <= index && index < len(path) {
			path = path[1+index:]
		}
	}

	return path
}
