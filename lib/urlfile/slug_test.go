package liburlfile_test

import (
	"testing"

	"github.com/reiver/socialfed-api/lib/urlfile"
)

func TestSlug(t *testing.T) {

	tests := []struct{
		Path string
		Expected string
	}{
		{

		},



		{
			Path:     "2025-02-07T12:20:57.323386567-08:00_M533NI1Z0RFI1PDTTODIJCQ35EYD00Y9L7I8M9B2.url",
			Expected: "2025-02-07T12:20:57.323386567-08:00",
		},



		{
			Path:     "/2025-02-07T12:20:57.323386567-08:00_M533NI1Z0RFI1PDTTODIJCQ35EYD00Y9L7I8M9B2.url",
			Expected:  "2025-02-07T12:20:57.323386567-08:00",
		},
		{
			Path:     "somefeed/2025-02-07T12:20:57.323386567-08:00_M533NI1Z0RFI1PDTTODIJCQ35EYD00Y9L7I8M9B2.url",
			Expected:          "2025-02-07T12:20:57.323386567-08:00",
		},
		{
			Path:     "/somefeed/2025-02-07T12:20:57.323386567-08:00_M533NI1Z0RFI1PDTTODIJCQ35EYD00Y9L7I8M9B2.url",
			Expected:           "2025-02-07T12:20:57.323386567-08:00",
		},
		{
			Path:     "feeds/somefeed/2025-02-07T12:20:57.323386567-08:00_M533NI1Z0RFI1PDTTODIJCQ35EYD00Y9L7I8M9B2.url",
			Expected:                "2025-02-07T12:20:57.323386567-08:00",
		},
	}

	for testNumber, test := range tests {

		actual := liburlfile.Slug(test.Path)

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual published date-time is not what was expected.", testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("PATH: %q", test.Path)
			continue
		}
	}
}
