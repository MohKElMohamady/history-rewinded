package cassandra

import (
	"crypto/tls"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	clientPool.Put(stargateClient1)

	stargateClient2 := clientPool.New().(*client.StargateClient)
	clientPool.Put(stargateClient2)

	stargateClient3 := clientPool.New().(*client.StargateClient)
	clientPool.Put(stargateClient3)
}
