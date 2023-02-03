package db

import (
	"context"
	"crypto/tls"
	"history-rewinded-goneril/models"
	"history-rewinded-regan/pb"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var remoteCassandraUri string
var remoteCassandraBearerToken string
var allocatedConnectionsToCassandra int64
var cassandraClientsPool = sync.Pool{
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
			log.Fatalf("failed to connect to the remote cassandra instance, reason: %v\n", err)
		}

		stargateClient, err := client.NewStargateClientWithConn(conn)
		if err != nil {
			log.Fatalf("failed to create instance of stargate client, reason: %v\n", err.Error())
		}

		log.Println("The total number of created instances is", allocatedConnectionsToCassandra)

		return stargateClient
	},
}

func init() {

	godotenv.Load()
	remoteCassandraUri = os.Getenv("CASSANDRA_REMOTE_URI")
	remoteCassandraBearerToken = os.Getenv("CASSANDRA_BEARER_TOKEN")

	client1 := cassandraClientsPool.Get().(client.StargateClient)
	defer cassandraClientsPool.Put(client1)
	go func() {
		client1.ExecuteQueryWithContext(
			&proto.Query{Cql: `CREATE TABLE IF NOT EXISTS goneril.users 
								(email text, username text, password text, birthday date) 
								PRIMARY KEY ((email, username))`},
			context.Background(),
		)
	}()

	client2 := cassandraClientsPool.Get().(client.StargateClient)
	defer cassandraClientsPool.Put(client2)
	go func() {
		client2.ExecuteQueryWithContext(
			&proto.Query{Cql: `CREATE TABLE IF NOT EXISTS goneril.favorites_incidents 
								(email text, day int, month int, summary text, incident_type text) 
								PRIMARY KEY ((email), day, month, summary)`},
			context.Background(),
		)
	}()

}

type CassandraRepository interface {
	RegisterUser(username, password string) error
	LoginWithEmail(email, password string) models.User
	LoginWithUsername(username, password string) models.User
	ResetPassword(email string) error
	AddFavoriteIncident(pb.Incident) error
	DeleteIncidentFromFavorites(id uuid.UUID) error
	FetchAllFavoriteIncidents() []pb.Incident
	FetchFavoriteEventsOn(day, month uint) []pb.Incident
	FetchFavoriteBirthsOn(day, month uint) []pb.Incident
	FetchFavoriteDeathsOn(day, month uint) []pb.Incident
	FetchFavoriteHolidaysOn(day, month uint) []pb.Incident
}

type CassandraClient struct {
}

func (c CassandraClient) RegisterUser(username, password string) error {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) LoginWithEmail(email, password string) models.User {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) LoginWithUsername(username, password string) models.User {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) ResetPassword(email string) error {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) AddFavoriteIncident(incident pb.Incident) error {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) DeleteIncidentFromFavorites(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) FetchAllFavoriteIncidents() []pb.Incident {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) FetchFavoriteEventsOn(day, month uint) []pb.Incident {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) FetchFavoriteBirthsOn(day, month uint) []pb.Incident {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) FetchFavoriteDeathsOn(day, month uint) []pb.Incident {
	//TODO implement me
	panic("implement me")
}

func (c CassandraClient) FetchFavoriteHolidaysOn(day, month uint) []pb.Incident {
	//TODO implement me
	panic("implement me")
}
