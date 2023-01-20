package wikipedia

import (
	"fmt"
	"history-rewinded-lear/models"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

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
		// log.Printf("%s", incident)
	}
	// return &models.Incident{Summary: r.Get("text").Str}
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
	case 244 < dayInYear && dayInYear < 274:
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
