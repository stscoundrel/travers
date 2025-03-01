package events

import (
	"log"

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
	log.Printf("Fetching events")
	freshEvents, err := s.GetAllEvents()
	if err != nil {
		return nil, err
	}
	log.Printf("%d total events found via fetchers.", len(freshEvents))

	log.Printf("Checking if events are new.")
	newEvents, err := s.repository.GetNewEvents(freshEvents)
	if err != nil {
		return nil, err
	}

	log.Printf("%d new events found.", len(newEvents))

	if len(newEvents) > 0 {
		log.Printf("Storing new events")
		err = s.repository.SaveEvents(newEvents)
		if err != nil {
			return nil, err
		}
	}

	return newEvents, nil
}
