package events

type EventService struct {
	sources []EventSource
}

func NewEventService(sources ...EventSource) *EventService {
	return &EventService{sources: sources}
}

func (s *EventService) GetAllEvents() ([]Event, error) {
	var allEvents []Event
	for _, source := range s.sources {
		events, err := source.FetchEvents()
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, events...)
	}
	return allEvents, nil
}
