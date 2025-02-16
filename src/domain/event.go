package domain

import "time"

type Event struct {
	ID        int
	Name      string
	Location  string
	StartTime time.Time
	EndTime   time.Time
}
