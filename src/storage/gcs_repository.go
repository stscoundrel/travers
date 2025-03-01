package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"github.com/stscoundrel/travers/domain"
)

type GCSRepository struct {
	bucketName string
	objectName string
	client     *storage.Client
}

func NewGCSRepository(bucketName, objectName string) (*GCSRepository, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSRepository{
		bucketName: bucketName,
		objectName: objectName,
		client:     client,
	}, nil
}

func (r *GCSRepository) SaveEvents(newEvents []domain.Event) error {
	storedEvents, err := r.GetEvents()
	if err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
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

	ctx := context.Background()
	writer := r.client.Bucket(r.bucketName).Object(r.objectName).NewWriter(ctx)
	defer writer.Close()

	_, err = writer.Write(data)
	return err
}

func (r *GCSRepository) GetEvents() ([]domain.Event, error) {
	ctx := context.Background()
	reader, err := r.client.Bucket(r.bucketName).Object(r.objectName).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			log.Printf("No existing events found in cloud storage")
			return []domain.Event{}, nil
		}
		return nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var storedEvents []domain.Event
	if err := json.Unmarshal(data, &storedEvents); err != nil {
		return nil, err
	}

	return storedEvents, nil
}

func (r *GCSRepository) GetNewEvents(freshEvents []domain.Event) ([]domain.Event, error) {
	storedEvents, err := r.GetEvents()
	if err != nil {
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
