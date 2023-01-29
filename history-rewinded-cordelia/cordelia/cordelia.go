package cordelia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"history-rewinded-cordelia/twitter"
	"history-rewinded-regan/pb"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
	"google.golang.org/grpc"
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
	http    *http.Client
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
	log.Println("successfully tweeted the event, the tweet id is %d and the its text is %s", tweet.TweetId, tweet.Text)
	return tweet 
}