package lear

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"history-rewinded-lear/cassandra"
	"history-rewinded-lear/models"
	"history-rewinded-lear/wikipedia"
	"history-rewinded-regan/pb"
	"log"
	"sync"
	"sync/atomic"
)

type lear struct {
	c cassandra.CassandraClient
	pb.UnimplementedLearServer
}

func New() lear {
	kingLear := lear{
		c: cassandra.CassandraClient{},
	}
	kingLear.FetchAllIncidents()
	return kingLear
}

func (l *lear) FetchAllIncidents() {
	incidentsSynchronizer := sync.WaitGroup{}
	eventsChannel := make(chan []*models.Incident)
	// birthsChannel := make(chan models.Incident)
	// deathsChannel := make(chan models.Incident)
	// holidaysChannel := make(chan models.Incident)

	incidentsSynchronizer.Add(1)
	go l.fetchAllEvents(eventsChannel)

	// incidentsSynchronizer.Add(1)
	//	go l.fetchAllBirths(birthsChannel)

	// incidentsSynchronizer.Add(1)
	//	go l.fetchAllDeaths(deathsChannel)

	// incidentsSynchronizer.Add(1)
	//	go l.fetchAllHolidays(holidaysChannel)

	// select {
	// case event := <- eventsChannel:
	// 	fmt.Printf("%s", event)
	// case birth := <- birthsChannel:
	// 	fmt.Printf("%s", birth)
	// case death := <- deathsChannel:
	// 	fmt.Printf("%s", death)
	// case holiday := <- holidaysChannel:
	// 	fmt.Printf("%s", holiday)
	// }

	go func() {
		for incidents := range eventsChannel {
			for _, _ = range incidents {
			}
			log.Println("Completed fetching all events and saved them")
		}
	}()

}

/*
 * Using the Fan-In pattern,
 */
func (l *lear) fetchAllEvents(fanInEventsChannel chan<- []*models.Incident) {

	// This WaitGroup is responsible for making sure that the 365 days of the years' events are fetched
	waitGroupEvents := sync.WaitGroup{}
	// One channel for every day in the year
	threeHundredSixtyFiveChannels := [365]chan []*models.Incident{}
	for i := 0; i < 365; i++ {
		threeHundredSixtyFiveChannels[i] = make(chan []*models.Incident)
	}
	// Add method on the WaitGroup must be outside spawning go routines to prevent race conditions
	waitGroupEvents.Add(len(threeHundredSixtyFiveChannels))
	// For each day of the year, we will send the responsible channel to fetch the data
	var totalEvents int64
	for i, incidentsChannel := range threeHundredSixtyFiveChannels {

		dayOfYear := i + 1

		go wikipedia.FetchIncidentFromWikipedia(incidentsChannel, models.Event, uint(dayOfYear))

		go func(incidentsChannel chan []*models.Incident) {
			defer waitGroupEvents.Done()
			for incidents := range incidentsChannel {
				for _, incident := range incidents {
					atomic.AddInt64(&totalEvents, 1)
					//It is faster to persist the incident right after parsing it, rather than waiting for it to be done by the fanning in channel
					l.c.AddIncident(*incident)
					log.Println(incident)
				}
				atomic.AddInt64(&totalEvents, 1)
				fanInEventsChannel <- incidents
			}
		}(incidentsChannel)

	}

	go func() {
		defer close(fanInEventsChannel)
		waitGroupEvents.Wait()
		//TODO: Calculate the hash out of all the events periodically every day and compare them with the database to know if the API has updated its information
		log.Printf("The total number of fetched events from the On This Day Wikipedia API is %d\n", totalEvents)
	}()

}

// func (l *lear) fetchAllBirths(fanInBirthsChannel chan <- models.Incident) {
// 	// This WaitGroup is responsible for making sure that the 365 days of the years' events are fetched
// 	waitGroupBirths := sync.WaitGroup{}
// 	// One channel for every day in the year
// 	threeHundredSixtyFiveChannels := make([] chan models.Incident, 365)
// 	// Add method on the WaitGroup must be outside spawning go routines to prevent race conditions
// 	waitGroupBirths.Add(365)

// 	// For each day of the year, we will send the responsible channel to fetch the data
// 	for dayOfYear, incidentCh := range threeHundredSixtyFiveChannels {

// 		defer waitGroupBirths.Done()

// 		go wikipedia.FetchIncidentFromWikipedia(incidentCh, models.Birth, rune(dayOfYear))

// 		go func(chan models.Incident){
// 			for incident := range incidentCh{
// 				fanInBirthsChannel <- incident
// 			}
// 		}(incidentCh)

// 	}

// 	go func(){
// 		waitGroupBirths.Wait()
// 		close(fanInBirthsChannel)
// 	}()
// }

// func (l *lear) fetchAllDeaths(deathsChannel chan <- models.Incident) {
// 	for dayOfYear := 1; dayOfYear <= 365; dayOfYear++ {
// 		go wikipedia.FetchIncidentFromWikipedia(models.Death, rune(dayOfYear))
// 	}
// 	close(deathsChannel)
// }

// func (l *lear) fetchAllHolidays(holidaysChannel chan <- models.Incident) {
// 	for dayOfYear := 1; dayOfYear <= 365; dayOfYear++ {
// 		go wikipedia.FetchIncidentFromWikipedia(models.Holidays, rune(dayOfYear))
// 	}
// 	close(holidaysChannel)
// }

func (l *lear) FetchIncidentsOn(*pb.FetchIncidentRequest, pb.Lear_FetchIncidentsOnServer) error {
	return status.Errorf(codes.Unimplemented, "method FetchIncidentsOn not implemented")
}

func (l *lear) FetchEventsOn(req *pb.FetchIncidentRequest, server pb.Lear_FetchEventsOnServer) error {

	events, err := l.c.FetchEventsOnThisDay(uint(req.Day), uint(req.Month))
	if err != nil {
		return err
	}

	for _, e := range events {
		server.Send(&pb.Incident{
			IncidentType:     pb.IncidentType_INCIDENT_TYPE_EVENT,
			Day:              e.Day,
			Month:            e.Month,
			Year:             e.Year,
			Summary:          e.Summary,
			IncidentInDetail: e.IncidentInDetail,
		})
	}

	return nil

}

func (l *lear) FetchBirthsOn(req *pb.FetchIncidentRequest, server pb.Lear_FetchBirthsOnServer) error {

	events, err := l.c.FetchBirthsOnThisDay(uint(req.Day), uint(req.Month))
	if err != nil {
		return err
	}

	for _, e := range events {
		server.Send(&pb.Incident{
			IncidentType:     pb.IncidentType_INCIDENT_TYPE_BIRTH,
			Day:              e.Day,
			Month:            e.Month,
			Year:             e.Year,
			Summary:          e.Summary,
			IncidentInDetail: e.IncidentInDetail,
		})
	}

	return nil

}
func (l *lear) FetchDeathsOn(req *pb.FetchIncidentRequest, server pb.Lear_FetchDeathsOnServer) error {

	events, err := l.c.FetchDeathsOnThisDay(uint(req.Day), uint(req.Month))
	if err != nil {
		return err
	}

	for _, e := range events {
		server.Send(&pb.Incident{
			IncidentType:     pb.IncidentType_INCIDENT_TYPE_DEATH,
			Day:              e.Day,
			Month:            e.Month,
			Year:             e.Year,
			Summary:          e.Summary,
			IncidentInDetail: e.IncidentInDetail,
		})
	}

	return nil

}
func (l *lear) FetchHolidaysOn(req *pb.FetchIncidentRequest, server pb.Lear_FetchHolidaysOnServer) error {

	events, err := l.c.FetchHolidaysOnThisDay(uint(req.Day), uint(req.Month))
	if err != nil {
		return err
	}

	for _, e := range events {

		server.Send(&pb.Incident{
			IncidentType:     pb.IncidentType_INCIDENT_TYPE_HOLIDAY,
			Day:              e.Day,
			Month:            e.Month,
			Year:             e.Year,
			Summary:          e.Summary,
			IncidentInDetail: e.IncidentInDetail,
		})
	}

	return nil

}
