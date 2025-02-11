package libfeeds

import (
	"testing"
)

func TestFeedPath(t *testing.T) {

	tests := []struct{
		Root string
		Name string
		Expected string
	}{
		{
			Root:     "/home/me/feed",
			Name:                   "something",
			Expected: "/home/me/feed/something",
		},
	}

	for testNumber, test := range tests {

		actual := feedPath(test.Root, test.Name)

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual path it not what was expected." , testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("ROOT: %q", test.Root)
			t.Logf("NAME: %q", test.Name)
			continue
		}
	}
}
