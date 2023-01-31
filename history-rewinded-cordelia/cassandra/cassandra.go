package cassandra

import (
	"crypto/tls"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	datastax "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"history-rewinded-cordelia/twitter"
	"log"
	"sync"
	"time"
)

var remoteCassandraUri string
var remoteCassandraBearerToken string
var clientPool = sync.Pool{
	New: func() any {

		config := &tls.Config{InsecureSkipVerify: false}
		conn, err := grpc.Dial(
			remoteCassandraUri,
			grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
			grpc.WithPerRPCCredentials(auth.NewStaticTokenProvider(remoteCassandraBearerToken)),
		)
		if err != nil {
			log.Fatalf("failed to connect to remote cassandra instance, reason:%s\n", err.Error())
		}

		stargateClient, err := client.NewStargateClientWithConn(conn)
		if err != nil {
			log.Fatalf("failed to instance instance of stargate client %s\n", err.Error())
		}

		return stargateClient
	},
}

func init() {

	client1 := clientPool.New().(*client.StargateClient)
	client2 := clientPool.New().(*client.StargateClient)
	client3 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(client1)
	defer clientPool.Put(client2)
	defer clientPool.Put(client3)

	go func() {
		client1.ExecuteQuery(&datastax.Query{
			Cql: `	CREATE TABLE IF NOT EXISTS cordelia.successful_tweets_by_id (tweeted_on int PRIMARY KEY, tweet_id int, tweet_text text);`,
		})
	}()
	
	go func() {
		client2.ExecuteQuery(&datastax.Query{
			Cql: `	CREATE TABLE IF NOT EXISTS cordelia.unsuccessful_tweets_by_timestamp 
					(attempted_to_tweet_on int PRIMARY KEY, reason text, status_code int , serialized_headers text);`,
		})
	}()

}

type CassandraRepository interface {
	SaveSuccessfulTweet(twitter.Tweet) error
	SaveUnsuccessfulTweet(twitter.TweetStatus) error
	FindTotalTweetsPerDay(time.Duration) uint64
}

type CassandraClient struct {
}

func (c *CassandraClient) SaveSuccessfulTweet(t twitter.Tweet) error {

	stargate := clientPool.Get().(*client.StargateClient)
	defer clientPool.Put(stargate)

	_, err := stargate.ExecuteQuery(&datastax.Query{
		Cql: "INSERT INTO cordelia.successful_tweets_by_id (tweeted_on, id, text) VALUES (?, ?, ?)",
		Values: &datastax.Values{
			Values: []*datastax.Value{
				{Inner: &datastax.Value_Int{Int: time.Now().Unix()}},
				{Inner: &datastax.Value_Int{Int: int64(t.TweetId)}},
				{Inner: &datastax.Value_String_{String_: t.Text}},
			}},
	})
	if err != nil {
		log.Printf("failed to save the tweet in the database, reason:%s\n", err.Error())
	}

	return nil

}

func (c *CassandraClient) SaveUnSuccessfulTweet(t twitter.TweetStatus) error {

	stargate := clientPool.Get().(*client.StargateClient)
	defer clientPool.Put(stargate)

	_, err := stargate.ExecuteQuery(&datastax.Query{
		Cql: "INSERT INTO cordelia.unsuccessful_tweets_by_timestamp (attempted_to_tweet_on, status_code, reason, serialized_headers) VALUES (?, ?, ?, ?)",
		Values: &datastax.Values{
			Values: []*datastax.Value{
				{Inner: &datastax.Value_Int{Int: time.Now().Unix()}},
				{Inner: &datastax.Value_Int{Int: int64(t.StatusCode)}},
				{Inner: &datastax.Value_String_{String_: t.Reason}},
				{Inner: &datastax.Value_String_{String_: t.SerializedHeaders}},
			}},
	})
	if err != nil {
		log.Printf("failed to save the unsuccessful tweet in the database, reason:%s\n", err.Error())
	}

	return nil

}