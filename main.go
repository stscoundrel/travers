package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stscoundrel/travers/domain"
	"github.com/stscoundrel/travers/events"
	"github.com/stscoundrel/travers/events/mpk"
	"github.com/stscoundrel/travers/storage"
)

func main() {
	// Local repository implementation.
	repo := storage.NewFileRepository("events.json")

	// Singular Santahamina data source.
	santahaminaFetcher := mpk.NewMPKEventSource(mpk.MPKQueryParams{
		UnitID:    mpk.UnitHelsinki,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0), // +1 year.
	})

	eventService := events.NewEventService(repo, santahaminaFetcher)

	// Fetch and store new events
	newEvents, err := eventService.FetchAndStoreNewEvents()
	if err != nil {
		log.Fatal("Error processing events:", err)
	}

	// Print results
	if len(newEvents) > 0 {
		fmt.Println("\n🔔 New events found:")
		printEvents(newEvents)
	} else {
		fmt.Println("\n📌 No new events found.")
	}

}

func printEvents(events []domain.Event) {
	for _, event := range events {
		fmt.Printf("Event: %s (%s) in %s\n",
			event.Name,
			event.StartTime.Format("02.01.2006 15.04"),
			event.Location)
	}
}
