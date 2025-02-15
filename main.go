package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stscoundrel/travers/events"
	"github.com/stscoundrel/travers/events/mpk"
)

func main() {
	santahaminaFetcher := mpk.NewMPKEventSource(mpk.MPKQueryParams{
		UnitID:    mpk.UnitHelsinki,
		StartDate: time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 9, 3, 0, 0, 0, 0, time.UTC),
	})

	upinniemiFetcher := mpk.NewMPKEventSource(mpk.MPKQueryParams{
		UnitID:    mpk.UnitUpinniemi,
		StartDate: time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 9, 3, 0, 0, 0, 0, time.UTC),
	})

	eventService := events.NewEventService(santahaminaFetcher, upinniemiFetcher)

	allEvents, err := eventService.GetAllEvents()
	if err != nil {
		log.Fatal("Error fetching events:", err)
	}

	for _, event := range allEvents {
		fmt.Printf("Event: %s (%s) in %s\n",
			event.Name,
			event.StartTime.Format("02.01.2006 15.04"),
			event.Location)
	}
}
