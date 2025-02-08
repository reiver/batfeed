package libfeeds

import (
	"fmt"
	"time"
)

// freshURLFileName returns a new name that can be used for a .url file.
//
// A .url file stores a single URL in it.
func freshURLFileName() string {
	now := time.Now()

	const layout string = time.RFC3339Nano
	nowstr := now.Format(layout)

	return fmt.Sprintf("%s_%s.url", nowstr, randstr())
}
