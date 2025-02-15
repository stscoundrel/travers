package mpk

import "time"

func parseTimeString(dateStr string) (time.Time, error) {
	layout := "2.1.2006 15.04"
	return time.Parse(layout, dateStr)
}
