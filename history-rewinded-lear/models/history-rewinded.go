package models

type IncidentType rune

// https://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-go
const (
	Event IncidentType = iota
	Birth
	Death
	Holidays
)

type Incident struct {
	Summary          string
	IncidentType     IncidentType
	IncidentInDetail string
	Year             int
}

