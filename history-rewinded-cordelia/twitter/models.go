package twitter

import (
	"fmt"
	"history-rewinded-regan/pb"
	"time"
)

type Tweet struct {
	TweetId uint64 `json:"tweet_id"`
	Text    string `json:"text"`
}

type TweetRequest struct {
	Text string `json:"text"`
}

type TweetStatus struct {
	
}

type IsTweetDeleted bool

func IncidentToTweetRequest(i *pb.Incident) Tweet {
	text := fmt.Sprintf("%v years ago, On this day %v-%v-%v, %s", (time.Now().Year() - int(i.Year)), i.Day, i.Month, i.Year, i.Summary)
	return Tweet{Text: text}	
}