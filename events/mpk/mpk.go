package mpk

import (
	"github.com/stscoundrel/travers/events"
)

type MPKEventSource struct {
	params MPKQueryParams
}

func NewMPKEventSource(params MPKQueryParams) *MPKEventSource {
	return &MPKEventSource{params: params}
}

func (m *MPKEventSource) FetchEvents() ([]events.Event, error) {
	return fetchEvents(m.params)
}
