package mpk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stscoundrel/travers/domain"
)

func fetchEvents(params MPKQueryParams) ([]domain.Event, error) {
	url := buildURL(params)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawEvents []rawMPKEvent
	if err := json.NewDecoder(resp.Body).Decode(&rawEvents); err != nil {
		return nil, err
	}

	return mapToDomainEvents(rawEvents)
}

func mapToDomainEvents(rawEvents []rawMPKEvent) ([]domain.Event, error) {
	var eventsList []domain.Event
	for _, raw := range rawEvents {
		startTime, _ := parseTimeString(raw.AlkuaikaStr)
		endTime, _ := parseTimeString(raw.LoppuaikaStr)

		eventsList = append(eventsList, domain.Event{
			ID:        raw.TapahtumaID,
			Name:      raw.Nimi,
			Location:  raw.Sijainti,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	return eventsList, nil
}

func buildURL(params MPKQueryParams) string {
	baseURL := "https://koulutuskalenteri.mpk.fi/Koulutuskalenteri"

	// Format start & end dates (DD.MM.YYYY)
	startDateStr := params.StartDate.Format("02.01.2006")
	endDateStr := params.EndDate.Format("02.01.2006")

	// Setup query parameters
	query := url.Values{}
	query.Set("type", "search")
	query.Set("format", "json")
	query.Set("unit_id", fmt.Sprintf("%d", params.UnitID))
	query.Set("start", startDateStr)
	query.Set("end", endDateStr)
	query.Set("only_my_events", "false")
	query.Set("VerkkoKoulutus", "false")
	query.Set("lisaysAikaleima", "false")
	query.Set("nayta_Vain_Ilmo_Auki", "false")
	query.Set("keyword_id", ShootingID)

	return fmt.Sprintf("%s?%s", baseURL, query.Encode())
}
