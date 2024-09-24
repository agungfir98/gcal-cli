package timeutils

import (
	"strings"
	"time"
)

var DefaultLayout = "02 01 2006"
var DefaultLayoutWithTime = "02 01 2006 15:04"
var AllDayDefaultLayout = "2006-01-02"

func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 23, 59, 59, 0, t.Location())
}

var layouts = []string{
	"02 01 2006 15:04",
	"02-01-2006 15:04",
	"02-01-2006 3:04 PM",
	"2006-01-02 15:04",
	"02/01/2006 15:04",
	"2006-01-02",
	"02/01/2006",
	"15:04",
	"3:04 PM",
}

var formatMap = map[string]string{
	"%a": "Mon",
	"%A": "Monday",
	"%b": "Jan",
	"%B": "January",
	"%d": "02",
	"%m": "01",
	"%Y": "2006",
	"%y": "06",
	"%H": "15",
	"%I": "03",
	"%p": "PM",
	"%M": "04",
	"%S": "05",
}

func ConvertToGoLayout(format string) string {
	for cFormat, goLayout := range formatMap {
		format = strings.Replace(format, cFormat, goLayout, -1)
	}
	return format
}

func ParseWithCustomFormat(format, input string) (time.Time, error) {
	t := time.Now()
	goLayout := ConvertToGoLayout(format)
	return time.ParseInLocation(goLayout, input, t.Location())
}
