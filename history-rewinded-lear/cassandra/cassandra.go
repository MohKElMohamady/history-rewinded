package cassandra

import (
	"crypto/tls"
	"fmt"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	datastax "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"history-rewinded-lear/models"
	"log"
	"os"
	"sync"
)

var remoteCassandraUri = os.Getenv("CASSANDRA_REMOTE_URI")
var remoteCassandraBearerToken = os.Getenv("CASSANDRA_BEARER_TOKEN")

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
			log.Fatalf("failed to connect to the remote cassandra instance, reason: %w", err)
			return nil
		}

		stargateClient, err := client.NewStargateClientWithConn(conn)

		return stargateClient
	},
}

// Initially, the pool will be having three clients so that it can be used, more will be created if needed
func init() {
	stargateClient1 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(stargateClient1)
	stargateClient2 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(stargateClient2)
	stargateClient3 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(stargateClient3)

	synchronizerChan := make(chan interface{})

	go func() {
		stargateClient1.ExecuteQuery(&datastax.Query{
			Cql: "CREATE KEYSPACE IF NOT EXISTS lear WITH replication = {'class': 'NetworkTopologyStrategy', 'europe-west1': '3'}  AND durable_writes = true;"},
		)
		close(synchronizerChan)
	}()

	// Design decision: Create four tables for
	go func() {
		<- synchronizerChan
		stargateClient1.ExecuteBatch(&datastax.Batch{
			Type: datastax.Batch_LOGGED,
			Queries: []*datastax.BatchQuery{
				{
					Cql: `CREATE TABLE IF NOT EXISTS lear.incidents (
					  incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
					  PRIMARY KEY (incident_id, year, month, day)
					  WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);
					`,
				},
				{
					Cql: `CREATE TABLE IF NOT EXISTS lear.births (
					  incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
					  PRIMARY KEY (incident_id, year, month, day)
					  WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);
					`,
				},
			},
		},
		)
	}()

	// Create the set of tables required the languages
   go func() {
	<- synchronizerChan	
	stargateClient2.ExecuteBatch(&datastax.Batch{
		Type: datastax.Batch_LOGGED,
		Queries: []*datastax.BatchQuery{
			{
				Cql: `CREATE TABLE IF NOT EXISTS lear.deaths (
					  incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
					  PRIMARY KEY (incident_id, year, month, day)
					  WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);
					`,
			},
			{
				Cql: `CREATE TABLE IF NOT EXISTS lear.holidays (
					  incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
					  PRIMARY KEY (incident_id, year, month, day)
					  WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);
					`,
			},
		},
	},
	)
   }()
}

