package events

import "time"

type Event struct {
	ID        int
	Name      string
	Location  string
	StartTime time.Time
	EndTime   time.Time
}

type EventSource interface {
	FetchEvents() ([]Event, error)
}
