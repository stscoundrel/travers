package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/stscoundrel/travers/domain"
)

type FileRepository struct {
	filePath string
}

func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{filePath: filePath}
}

func (r *FileRepository) SaveEvents(newEvents []domain.Event) error {
	storedEvents, err := r.GetEvents()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	existingEventIDs := make(map[int]bool)
	for _, event := range storedEvents {
		existingEventIDs[event.ID] = true
	}

	for _, event := range newEvents {
		if !existingEventIDs[event.ID] {
			storedEvents = append(storedEvents, event)
		}
	}

	data, err := json.MarshalIndent(storedEvents, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *FileRepository) GetEvents() ([]domain.Event, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("Failed to read local events DB")
		}
		return nil, err
	}

	var storedEvents []domain.Event
	if err := json.Unmarshal(data, &storedEvents); err != nil {
		return nil, err
	}

	return storedEvents, nil
}

func (r *FileRepository) GetNewEvents(freshEvents []domain.Event) ([]domain.Event, error) {
	storedEvents, err := r.GetEvents()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	existingEventIDs := make(map[int]bool)
	for _, event := range storedEvents {
		existingEventIDs[event.ID] = true
	}

	var newEvents []domain.Event
	for _, event := range freshEvents {
		if !existingEventIDs[event.ID] {
			newEvents = append(newEvents, event)
		}
	}

	return newEvents, nil
}
