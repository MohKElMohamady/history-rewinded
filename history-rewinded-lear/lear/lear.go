package lear

import (
	"history-rewinded-lear/models"
	"history-rewinded-lear/wikipedia"
	"log"
	"sync"
	"sync/atomic"
)

type lear struct {
}

func New() *lear {
	return &lear{}
}

func (l *lear) FetchAllIncidents() {

	eventsChannel := make(chan []*models.Incident)
	// birthsChannel := make(chan models.Incident)
	// deathsChannel := make(chan models.Incident)
	// holidaysChannel := make(chan models.Incident)

	go l.fetchAllEvents(eventsChannel)
	//	go l.fetchAllBirths(birthsChannel)
	//	go l.fetchAllDeaths(deathsChannel)
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
	for incidents := range eventsChannel {
		for _ = range incidents {
			// log.Println("Something")
			// log.Printf("from channel printer %v", incident)
		}
	}
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

			log.Println(incidentsChannel)

			go func(incidentsChannel * chan []*models.Incident) {
				for incidents := range *incidentsChannel {
					log.Printf("%v", len(incidents))
					for _, incident := range incidents {
						atomic.AddInt64(&totalEvents, 1)
						log.Println(incident)
					}
					fanInEventsChannel <- incidents
				}
				defer waitGroupEvents.Done()
			}(&incidentsChannel)

	}

		log.Println(totalEvents)
	go func() {
		waitGroupEvents.Wait()
		close(fanInEventsChannel)
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