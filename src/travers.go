package travers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	// Hook up the cloud function.
	functions.HTTP("FetchAndStoreEvents", FetchAndStoreEvents)

	log.Printf("Initialized HTTP endpoint")
}

// Cloud function entry point.
func FetchAndStoreEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("Travers entry point")

	bucketName := os.Getenv("EVENT_STORAGE_BUCKET")

	if bucketName == "" {
		log.Fatal("EVENT_STORAGE_BUCKET environment variable is not set!")
	}

	log.Printf("Using Cloud Storage bucket: %s", bucketName)

	repo, err := storage.NewGCSRepository(bucketName, "events.json")
	if err != nil {
		log.Fatalf("Failed to initialize GCS Repository: %v", err)
	}

	// Santahamina data source.
	santahaminaFetcher := mpk.NewMPKEventSource(mpk.MPKQueryParams{
		UnitID:    mpk.UnitHelsinki,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0), // +1 year.
	})

	log.Printf("Initialized Santahamina data source.")

	eventService = events.NewEventService(repo, santahaminaFetcher)

	log.Printf("Initialized Event service.")

	log.Printf("Fetching and storing new events")
	newEvents, err := eventService.FetchAndStoreNewEvents()
	if err != nil {
		log.Fatal("Error processing events:", err)
	}

	response := struct {
		Message string         `json:"message"`
		Events  []domain.Event `json:"events,omitempty"`
	}{
		Message: "No new events found.",
	}

	if len(newEvents) > 0 {
		log.Printf("%d new events found.", len(newEvents))
		response.Message = fmt.Sprintf("%d new events found.", len(newEvents))
		response.Events = newEvents
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
