package models

import "fmt"

type IncidentType rune

// https://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-go
const (
	Event IncidentType = iota
	Birth
	Death
	Holiday
)

func (iT IncidentType) String() string {
	switch iT {
	case Event:
		return "event"
	case Birth:
		return "birth"
	case Death:
		return "death"
	case Holiday:
		return "holiday"
	default:
		return ""
	}
}

type Incident struct {
	Summary          string
	IncidentType     IncidentType
	IncidentInDetail string
	Day              int64
	Month            int64
	Year             int64
}

func (i *Incident) String() string {
	return fmt.Sprintf("On this day, %d-%d-%d, %s\n", i.Day, i.Month, i.Year, i.Summary)
}
