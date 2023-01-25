package wikipedia

import (
	"fmt"
	"history-rewinded-lear/models"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"github.com/tidwall/gjson"
)


// https://stackoverflow.com/questions/36773837/best-way-to-use-http-client-in-a-concurrent-application
var httpClientsPool = sync.Pool{
	New: func() any { return &http.Client{} }}

const baseUrl = "https://en.wikipedia.org/api/rest_v1/feed/onthisday/"

func init() {
	httpClientsPool.Put(&http.Client{})
	httpClientsPool.Put(&http.Client{})
	httpClientsPool.Put(&http.Client{})
	httpClientsPool.Put(&http.Client{})
	httpClientsPool.Put(&http.Client{})
}

func FetchIncidentFromWikipedia(incidentChannel chan <- []*models.Incident, incidentType models.IncidentType, dayOfYear uint) {

	// https://stackoverflow.com/questions/62018471/why-does-golang-strings-builder-implements-string-like-this
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(baseUrl)

	// Get the day and the month
	month, day, err := GetDayMonthFromYearDay(dayOfYear)

	if err != nil {
		log.Printf("failed to convert the day of the year into month and day reason: %s", err.Error())
	}

	// What is the type of the event?
	switch incidentType {
	case models.Event:
		urlBuilder.WriteString("events")
	case models.Birth:
		urlBuilder.WriteString("births")
	case models.Death:
		urlBuilder.WriteString("deaths")
	case models.Holidays:
		urlBuilder.WriteString("holidays")
	}

	urlBuilder.WriteString(fmt.Sprintf("/%d/%d", month, day))

	// https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go
	httpClient := http.Client{}

	wikipediaResponse, err := httpClient.Get(urlBuilder.String())
	if err != nil {
		log.Printf("failed to fetch response from Wikipedia API, reason: %s", err)
	}
	defer wikipediaResponse.Body.Close()

	wikipediaBytesResponse, err := io.ReadAll(wikipediaResponse.Body)
	if err != nil {
		log.Printf("failed to parse the response into bytes, reason %v\n", err.Error())
	}

	parsedIncidents := ParseWikipediaResponseToIncidents(string(wikipediaBytesResponse), incidentType)

	incidentChannel <- parsedIncidents

	log.Println("Finished writing on channel, exiting function")
}

func ParseWikipediaResponseToIncidents(jsonBlob string, incidentType models.IncidentType) []*models.Incident {

	var incidents []*models.Incident

	// https://mholt.github.io/json-to-go/
	parsedJsonIncidents := gjson.Parse(jsonBlob).Array()[0].Get("events").Array()
	for _, r := range parsedJsonIncidents {

		incidentSummary := r.Get("text").Str
		incidentInDetail := r.Get("year").Str
		incidentYear := r.Get("year").Int()

		incident := &models.Incident{
			Summary:          incidentSummary,
			IncidentType:     incidentType,
			IncidentInDetail: incidentInDetail,
			Year:             int(incidentYear),
		}

		incidents = append(incidents, incident)
	}
	return incidents
}

var cumulativeDaysAtEndOfMonths = [12]uint{
	1, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334,
}

func GetDayMonthFromYearDay(dayInYear uint) (month uint, day uint, err error) {

	switch {
		
	case 1 <= dayInYear && dayInYear < 32:
		return 1, dayInYear, nil
	case 32 <= dayInYear && dayInYear < 60:
		return 2, (dayInYear - cumulativeDaysAtEndOfMonths[1]), nil
	case 60 <= dayInYear && dayInYear < 91:
		return 3, (dayInYear - cumulativeDaysAtEndOfMonths[2]), nil
	case 91 <= dayInYear && dayInYear < 121:
		return 4, (dayInYear - cumulativeDaysAtEndOfMonths[3]), nil
	case 121 <= dayInYear && dayInYear < 152:
		return 5, (dayInYear - cumulativeDaysAtEndOfMonths[4]), nil
	case 152 <= dayInYear && dayInYear < 182:
		return 6, (dayInYear - cumulativeDaysAtEndOfMonths[5]), nil
	case 182 <= dayInYear && dayInYear < 213:
		return 7, (dayInYear - cumulativeDaysAtEndOfMonths[6]), nil
	case 213 <= dayInYear && dayInYear < 244:
		return 8, (dayInYear - cumulativeDaysAtEndOfMonths[7]), nil
	case 244 <= dayInYear && dayInYear < 274:
		return 9, (dayInYear - cumulativeDaysAtEndOfMonths[8]), nil
	case 274 <= dayInYear && dayInYear < 305:
		return 10, (dayInYear - cumulativeDaysAtEndOfMonths[9]), nil
	case 305 <= dayInYear && dayInYear < 334:
		return 11, (dayInYear - cumulativeDaysAtEndOfMonths[10]), nil
	case 334 <= dayInYear && dayInYear <= 365:
		return 12, (dayInYear - cumulativeDaysAtEndOfMonths[11]), nil
	default:
		return 0, 0, fmt.Errorf("failed to parsed the month and day from %d", dayInYear)
	}
}

func GetDayOfYearFromDayAndMonth(month uint, day uint) (dayOfYear uint, err error) {
	switch {
	case month == 1:
		return day, nil
	case month == 2:
		return day + cumulativeDaysAtEndOfMonths[1], nil
	case month == 3:
		return day + cumulativeDaysAtEndOfMonths[2], nil
	case month == 4:
		return day + cumulativeDaysAtEndOfMonths[3], nil
	case month == 5:
		return day + cumulativeDaysAtEndOfMonths[4], nil
	case month == 6:
		return day + cumulativeDaysAtEndOfMonths[5], nil
	case month == 7:
		return day + cumulativeDaysAtEndOfMonths[6], nil
	case month == 8:
		return day + cumulativeDaysAtEndOfMonths[7], nil
	case month == 9:
		return day + cumulativeDaysAtEndOfMonths[8], nil
	case month == 10:
		return day + cumulativeDaysAtEndOfMonths[9], nil
	case month == 11:
		return day + cumulativeDaysAtEndOfMonths[10], nil
	case month == 12:
		return day + cumulativeDaysAtEndOfMonths[11], nil
	default: 
		return 0, fmt.Errorf("failed to get the day of the year for the month %v and day %v", month, day)
	} 
}