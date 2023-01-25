package models

import "fmt"

type IncidentType rune

// https://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-go
const (
	Event IncidentType = iota
	Birth
	Death
	Holidays
)

func (iT IncidentType) String() string {
	switch iT {
	case Event:
		return "Event"	
	case Birth:
		return "Birth"
	case Death:
		return "Death"
	case Holidays:
		return "Holidays"
	default:
		return ""
	}
}

type Incident struct {
	Summary          string
	IncidentType     IncidentType
	IncidentInDetail string
	Day int
	Month int
	Year             int
}

func (i *Incident) String() string {
	return fmt.Sprintf("On this day, %d-%d-%d, %s\n", i.Day, i.Month, i.Year, i.Summary)	
}