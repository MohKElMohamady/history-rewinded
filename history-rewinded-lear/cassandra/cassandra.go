package cassandra

import (
	"crypto/tls"
	"errors"
	"fmt"
	"history-rewinded-lear/models"
	"log"
	"os"
	"sync"
	"time"
	"github.com/joho/godotenv"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	datastax "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
			log.Fatalf("failed to connect to the remote cassandra instance, reason: %w", err)
			return nil
		}

		stargateClient, err := client.NewStargateClientWithConn(conn)

		return stargateClient
	},
}

// Initially, the pool will be having three clients so that it can be used, more will be created if needed
func init() {
	godotenv.Load()
	remoteCassandraUri = os.Getenv("CASSANDRA_REMOTE_URI")
	remoteCassandraBearerToken = os.Getenv("CASSANDRA_BEARER_TOKEN")
	stargateClient1 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(stargateClient1)
	stargateClient2 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(stargateClient2)
	stargateClient3 := clientPool.New().(*client.StargateClient)
	defer clientPool.Put(stargateClient3)

	synchronizerChan := make(chan interface{})

	go func() {
		// _, err :=	stargateClient1.ExecuteQuery(&datastax.Query{
		// 		Cql: "CREATE KEYSPACE IF NOT EXISTS lear WITH replication = {'class': 'NetworkTopologyStrategy', 'europe-west1': '3'}  AND durable_writes = true;"},
		// 	)
		// if err != nil {
		// 	log.Fatalln(err.Error())
		// }
		close(synchronizerChan)
	}()

	// Design decision: Create four tables for each incident type
	go func() {
		<-synchronizerChan
		_, err := stargateClient2.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.incidents (
				incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
				PRIMARY KEY (incident_id, year, month, day))
				WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = stargateClient2.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.births (
				incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
				PRIMARY KEY (incident_id, year, month, day))
				WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	go func() {
		<-synchronizerChan
		_, err := stargateClient3.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.deaths (
				incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
				PRIMARY KEY (incident_id, year, month, day))
				WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = stargateClient3.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.holidays (
				incident_id uuid, summary text, incident_in_detail text, year int, month int, day int,
				PRIMARY KEY (incident_id, year, month, day))
				WITH CLUSTERING ORDER BY (year DESC, month ASC, day ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()
}

type CassandraRepository interface {
	FetchIncidentsOnThisDay(day uint, month uint) ([]models.Incident, error)
	FetchBirthsOnThisDay(day uint, month uint) ([]models.Incident, error)
	FetchDeathsOnThisDay(day uint, month uint) ([]models.Incident, error)
	FetchHolidaysOnThisDay(day uint, month uint) ([]models.Incident, error)
	AddIncident(incident models.Incident) (models.Incident, error)
	BulkInsertIncidents(incidents models.Incident) ([]models.Incident, error)
}

type CassandraClient struct {
}

func (c *CassandraClient) FetchIncidentsOnThisDay(day uint, month uint) ([]models.Incident, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CassandraClient) FetchBirthsOnThisDay(day uint, month uint) ([]models.Incident, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CassandraClient) FetchDeathsOnThisDay(day uint, month uint) ([]models.Incident, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CassandraClient) FetchHolidaysOnThisDay(day uint, month uint) ([]models.Incident, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CassandraClient) AddIncident(incident models.Incident) (models.Incident, error) {

	stargate := clientPool.Get().(*client.StargateClient)
	switch {
	case time.Now().Year()<= incident.Year :
		return nil, fmt.Errorf("the event cannot be in the future, the year value is %v", incident.Month)
	case incident.Month <= 0 || incident.Month > 12:
		return nil, fmt.Errorf("the month cannot be less than 1 as it has the value %v", incident.Month)
	case incident.Day <= 0 || incident.Month > 31:
		return nil, fmt.Errorf("the day cannot be less than 1 as it has the value %v", incident.Day)
	case incident.Summary == "":
		return nil, errors.New("the incident summary cannot be empty")
	case incident.IncidentInDetail == "":
		return nil, errors.New("the incident summary cannot be empty")
	case incident.IncidentType != models.Event || incident.IncidentType != models.Birth ||incident.IncidentType != models.Death || incident.IncidentType != models.Holiday :
		return nil, errors.New("unknown incident type")
	default:
		break
	}

	_, err := stargate.ExecuteQuery(&datastax.Query{
		Cql: "INSERT INTO lear.incidents (incident_id , year , month , day , summary , incident_in_detail ) VALUES ( uuid(), ?, ?, ?, ?, ?);",
		Values: &datastax.Values{Values: []*datastax.Value{
			&datastax.Value{Inner: &datastax.Value_Int{Int: int64(incident.Year)}},
			&datastax.Value{Inner: &datastax.Value_Int{int64(incident.Month)}},
			&datastax.Value{Inner: &datastax.Value_Int{int64(incident.Day)}},
			&datastax.Value{Inner: &datastax.Value_String_{String_:incident.Summary}},
			&datastax.Value{Inner: &datastax.Value_String_{String_:incident.IncidentInDetail}},
		}},
	})
	if err != nil {
		log.Printf("error, unable to insert incident into its table, reason: %v ", err.Error())
	}

	return incident, nil
}

func (c *CassandraClient) BulkInsertIncidents(incidents models.Incident) ([]models.Incident, error) {
	//TODO implement me
	panic("implement me")
}