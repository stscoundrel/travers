package mpk

import "github.com/stscoundrel/travers/domain"

type MPKEventSource struct {
	params MPKQueryParams
}

func NewMPKEventSource(params MPKQueryParams) *MPKEventSource {
	return &MPKEventSource{params: params}
}

func (m *MPKEventSource) FetchEvents() ([]domain.Event, error) {
	return fetchEvents(m.params)
}
