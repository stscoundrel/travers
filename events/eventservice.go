package events

import (
	"github.com/stscoundrel/travers/domain"
	"github.com/stscoundrel/travers/storage"
)

type EventSource interface {
	FetchEvents() ([]domain.Event, error)
}

type EventService struct {
	sources    []EventSource
	repository storage.Repository
}

func NewEventService(repo storage.Repository, sources ...EventSource) *EventService {
	return &EventService{
		sources:    sources,
		repository: repo,
	}
}

func (s *EventService) GetAllEvents() ([]domain.Event, error) {
	var allEvents []domain.Event
	for _, source := range s.sources {
		events, err := source.FetchEvents()
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, events...)
	}
	return allEvents, nil
}

func (s *EventService) FetchAndStoreNewEvents() ([]domain.Event, error) {
	freshEvents, err := s.GetAllEvents()
	if err != nil {
		return nil, err
	}

	// Identify new events
	newEvents, err := s.repository.GetNewEvents(freshEvents)
	if err != nil {
		return nil, err
	}

	// Store only new events
	if len(newEvents) > 0 {
		err = s.repository.SaveEvents(newEvents)
		if err != nil {
			return nil, err
		}
	}

	return newEvents, nil
}
