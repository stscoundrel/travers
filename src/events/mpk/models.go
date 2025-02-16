package mpk

import "time"

// UnitID represents an MPK training unit.
type UnitID int

const (
	UnitHelsinki  UnitID = 2
	UnitUpinniemi UnitID = 21
)

const ShootingID = "22"

func (u UnitID) String() string {
	switch u {
	case UnitHelsinki:
		return "Helsinki"
	case UnitUpinniemi:
		return "Upinniemi"
	default:
		return "Unknown"
	}
}

type MPKQueryParams struct {
	UnitID    UnitID
	StartDate time.Time
	EndDate   time.Time
}

type rawMPKEvent struct {
	TapahtumaID  int    `json:"TapahtumaID"`
	Nimi         string `json:"Nimi"`
	Sijainti     string `json:"PostitoimipaikkaListassa"`
	AlkuaikaStr  string `json:"AlkuaikaStr"`
	LoppuaikaStr string `json:"LoppuaikaStr"`
}
