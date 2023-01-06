package cordelia

import (
	"context"
	"fmt"
	"github.com/dghubble/oauth1"
	"net/http"
	"os"
)

type cordelia struct {
	client *http.Client
}

func New() (cordelia, error) {
	config := oauth1.NewConfig(os.Getenv("TWITTER_OAUTH_1_API_KEY"), os.Getenv("TWITTER_OAUTH_1_API_KEY_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_OAUTH_1_ACCESS_TOKEN"), os.Getenv("TWITTER_OAUTH_1_ACCESS_TOKEN_SECRET"))
	httpClient := oauth1.NewClient(context.Background(), config, token)
	return cordelia{httpClient}, nil
}

func (c *cordelia) Tweet() {
	fmt.Println("Tweeting, tweet tweet")
}
