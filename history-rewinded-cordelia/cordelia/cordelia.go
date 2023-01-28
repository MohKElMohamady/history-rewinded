package cordelia

import (
	"context"
	"fmt"
	"history-rewinded-cordelia/twitter"
	"history-rewinded-regan/pb"
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
	client    *http.Client
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
