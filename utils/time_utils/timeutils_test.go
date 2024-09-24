package timeutils

import (
	"testing"
)

func TestConverToGoLayout(t *testing.T) {

	testCase := []struct {
		title, value, expected string
	}{
		{"%d-%m-%Y to 01-02-2006", "%d-%m-%Y", "02-01-2006"},
		{"%d-%m-%Y %H:%M to 01-02-2006 15:04", "%d-%m-%Y %H:%M", "02-01-2006 15:04"},
		{"%A %d %B %Y to Monday 02 01 2006", "%A %d %B %Y", "Monday 02 January 2006"},
	}

	for _, tc := range testCase {
		t.Run(tc.title, func(t *testing.T) {
			result := ConvertToGoLayout(tc.value)
			if result != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, result)
			}
		})
	}

}
