package storage

import (
	"github.com/stscoundrel/travers/domain"
)

type Repository interface {
	SaveEvents(events []domain.Event) error
	GetEvents() ([]domain.Event, error)
	GetNewEvents(freshEvents []domain.Event) ([]domain.Event, error)
}
