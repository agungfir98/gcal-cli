package timeutils

import "time"

var DefaultLayout = "02 01 2006 15:04"

func EndOfDay() time.Time {
	t := time.Now()
	y, m, d := t.Date()

	return time.Date(y, m, d, 23, 59, 59, 0, t.Location())
}
