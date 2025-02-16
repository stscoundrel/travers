package travers

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/stscoundrel/travers/domain"
	"github.com/stscoundrel/travers/events"
	"github.com/stscoundrel/travers/events/mpk"
	"github.com/stscoundrel/travers/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var eventService *events.EventService

func init() {
	// Local repository implementation.
	// TODO: cloud storage.
	repo := storage.NewFileRepository("events.json")

	// Santahamina data source.
	santahaminaFetcher := mpk.NewMPKEventSource(mpk.MPKQueryParams{
		UnitID:    mpk.UnitHelsinki,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0), // +1 year.
	})

	eventService = events.NewEventService(repo, santahaminaFetcher)

	// Hook up the cloud function.
	functions.HTTP("FetchAndStoreEvents", FetchAndStoreEvents)
}

// Cloud function entry point.
func FetchAndStoreEvents(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprint(w, "Irrelevant response")

}

func printEvents(events []domain.Event) {
	for _, event := range events {
		fmt.Printf("Event: %s (%s) in %s\n",
			event.Name,
			event.StartTime.Format("02.01.2006 15.04"),
			event.Location)
	}
}
