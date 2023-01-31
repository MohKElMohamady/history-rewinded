package cordelia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"google.golang.org/grpc"
	"history-rewinded-cordelia/twitter"
	"history-rewinded-regan/pb"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const kingLearAddress = ":1606"

type Tweeter interface {
	TweetEvent(*pb.Incident) twitter.Tweet
	TweetBirth(*pb.Incident) twitter.Tweet
	TweetDeath(*pb.Incident) twitter.Tweet
	TweetHoliday(*pb.Incident) twitter.Tweet
	DeleteTweet(tweetId uint64) twitter.IsTweetDeleted
	ReplyToTweetWithMoreInformation(furtherInformation string) twitter.Tweet
}

type cordelia struct {
	http      *http.Client
	king_lear pb.LearClient
}

func New() (cordelia, error) {

	conn, err := grpc.Dial(kingLearAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to establish connection with king lear, reason: %s", conn)
	}

	king_lear := pb.NewLearClient(conn)

	config := oauth1.NewConfig(os.Getenv("TWITTER_OAUTH_1_API_KEY"), os.Getenv("TWITTER_OAUTH_1_API_KEY_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_OAUTH_1_ACCESS_TOKEN"), os.Getenv("TWITTER_OAUTH_1_ACCESS_TOKEN_SECRET"))
	httpClient := oauth1.NewClient(context.Background(), config, token)
	return cordelia{httpClient, king_lear}, nil

}

func (c *cordelia) Tweet() {
	fmt.Println("Tweeting, tweet tweet")
}

func (c *cordelia) TweetEvent(incident *pb.Incident) twitter.Tweet {
	tweetReqest := twitter.IncidentToTweetRequest(incident)
	marshalledTweetRequest, err := json.Marshal(tweetReqest)
	if err != nil {
		log.Fatalf("")
	}

	res, err := c.http.Post("https://api.twitter.com/2/tweets", "application/json", bytes.NewBuffer(marshalledTweetRequest))
	if err != nil {
		//TODO: Save the error with its status and reason
		log.Fatalf("failed to tweet the event, reason:%v\n", err.Error())
	}

	defer res.Body.Close()
	var tweet twitter.Tweet

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &tweet)
	if err != nil {
		log.Fatalf("failed to unmarshalled the tweet from Twitter api, reason:%s\n", err.Error())
	}

	//TODO: Save the successfull tweet with its timestamp in Cassandra
	log.Printf("successfully tweeted the event, the tweet id is %d and the its text is %s", tweet.TweetId, tweet.Text)
	return tweet
}

func (c *cordelia) FetchIncidentaOfTheDay() (events []pb.Incident, births []pb.Incident, deaths []pb.Incident, holidays []pb.Incident) {

	today := int64(time.Now().Day())
	currentMonth := int64(time.Now().Month())

	fourIncidentsTypeSynchronizer := sync.WaitGroup{}

	fourIncidentsTypeSynchronizer.Add(4)
	events = []pb.Incident{}
	births = []pb.Incident{}
	deaths = []pb.Incident{}
	holidays = []pb.Incident{}

	eventsChannel := c.FetchEventsOn(today, currentMonth, pb.IncidentType_INCIDENT_TYPE_EVENT)
	birthsChannel := c.FetchEventsOn(today, currentMonth, pb.IncidentType_INCIDENT_TYPE_EVENT)
	deathsChannel := c.FetchEventsOn(today, currentMonth, pb.IncidentType_INCIDENT_TYPE_EVENT)
	holidaysChannel := c.FetchEventsOn(today, currentMonth, pb.IncidentType_INCIDENT_TYPE_EVENT)

	go func() {
		for e := range eventsChannel {
			events = append(events, e)
		}
		fourIncidentsTypeSynchronizer.Done()
	}()

	go func() {
		for b := range birthsChannel {
			events = append(events, b)
		}
		fourIncidentsTypeSynchronizer.Done()
	}()

	go func() {
		for d := range deathsChannel {
			events = append(events, d)
		}
		fourIncidentsTypeSynchronizer.Done()
	}()

	go func() {
		for h := range holidaysChannel {
			events = append(events, h)
		}
		fourIncidentsTypeSynchronizer.Done()
	}()

	fourIncidentsTypeSynchronizer.Wait()
	return events, births, deaths, holidays
}

func (c *cordelia) FetchEventsOn(day int64, month int64, incidentType pb.IncidentType) <-chan pb.Incident {

	EventsIncidents := make(chan pb.Incident)

	req := pb.FetchIncidentRequest{
		Day:   day,
		Month: month,
	}
	go func() {
		switch incidentType {
		case pb.IncidentType_INCIDENT_TYPE_EVENT:
			client, err := c.king_lear.FetchEventsOn(context.Background(), &req)
			log.Fatalf("failed to invoke the rpc for fetching events, reason:%s\n", err.Error())
			for {
				incident, err := client.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("failed to continue receive incidents from stream, reason:%s\n", err.Error())
				}
				EventsIncidents <- *incident
			}
		case pb.IncidentType_INCIDENT_TYPE_BIRTH:
			client, err := c.king_lear.FetchBirthsOn(context.Background(), &req)
			log.Fatalf("failed to invoke the rpc for fetching births, reason:%s\n", err.Error())
			for {
				incident, err := client.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("failed to continue receive births from stream, reason:%s\n", err.Error())
				}
				EventsIncidents <- *incident
			}
		case pb.IncidentType_INCIDENT_TYPE_DEATH:
			client, err := c.king_lear.FetchEventsOn(context.Background(), &req)
			log.Fatalf("failed to invoke the rpc for fetching deaths, reason:%s\n", err.Error())
			for {
				incident, err := client.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("failed to continue receive deaths from stream, reason:%s\n", err.Error())
				}
				EventsIncidents <- *incident
			}
		case pb.IncidentType_INCIDENT_TYPE_HOLIDAY:
			client, err := c.king_lear.FetchEventsOn(context.Background(), &req)
			log.Fatalf("failed to invoke the rpc for fetching holidays, reason:%s\n", err.Error())
			for {
				incident, err := client.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("failed to continue receive holidays from stream, reason:%s\n", err.Error())
				}
				EventsIncidents <- *incident
			}

		}
	}()

	return EventsIncidents
}
