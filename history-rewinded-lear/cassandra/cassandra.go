package cassandra

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	datastax "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"history-rewinded-lear/models"
	"log"
	"os"
	"sync"
	"sync/atomic"
	// "time"
	// "errors"
)

var remoteCassandraUri string
var remoteCassandraBearerToken string
var allocatedConnectionsToCassandra int64
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
		atomic.AddInt64(&allocatedConnectionsToCassandra, 1)
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
			Cql: `CREATE TABLE IF NOT EXISTS lear.events (
				day int, month int, year int, summary text, incident_in_detail text,
				PRIMARY KEY ((day, month), year, summary))
				WITH CLUSTERING ORDER BY (year DESC, summary ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = stargateClient2.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.deaths (
				day int, month int, year int, summary text, incident_in_detail text,
				PRIMARY KEY ((day, month), year, summary))
				WITH CLUSTERING ORDER BY (year DESC, summary ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	go func() {
		<-synchronizerChan
		_, err := stargateClient3.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.births (
				day int, month int, year int, summary text, incident_in_detail text,
				PRIMARY KEY ((month, day), year, summary))
				WITH CLUSTERING ORDER BY (year DESC, summary ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}

		_, err = stargateClient3.ExecuteQuery(&datastax.Query{
			Cql: `CREATE TABLE IF NOT EXISTS lear.holidays (
				day int, month int, year int, summary text, incident_in_detail text,
				PRIMARY KEY ((month, day), year, summary))
				WITH CLUSTERING ORDER BY (year DESC, summary ASC);`,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()
}

type CassandraRepository interface {
	FetchEventsOnThisDay(day uint, month uint) ([]models.Incident, error)
	FetchBirthsOnThisDay(day uint, month uint) ([]models.Incident, error)
	FetchDeathsOnThisDay(day uint, month uint) ([]models.Incident, error)
	FetchHolidaysOnThisDay(day uint, month uint) ([]models.Incident, error)
	AddIncident(incident models.Incident) (models.Incident, error)
	BulkInsertIncidents(incidents models.Incident) ([]models.Incident, error)
}

type CassandraClient struct {
}

func (c *CassandraClient) FetchEventsOnThisDay(day uint, month uint) ([]models.Incident, error) {
	stargateClient := clientPool.Get().(*client.StargateClient)
	defer clientPool.Put(stargateClient)
	fetchIncidentQuery := &datastax.Query{
		Cql: "SELECT year, summary, incident_in_summary FROM lear.incidents WHERE day = ? AND month = ?",
		Values: &datastax.Values{
			Values: []*datastax.Value{
				{
					Inner: &datastax.Value_Int{Int: int64(day)},
				},
				{
					Inner: &datastax.Value_Int{Int: int64(month)},
				},
			},
		},
	}

	resp, err := stargateClient.ExecuteQuery(fetchIncidentQuery)
	if err != nil {
		fmt.Printf("failed to fetch results from remote cassandra instance, reason %s\n", err.Error())
	}

	incidents := []models.Incident{}

	resultSet := resp.GetResultSet()

	for _, row := range resultSet.Rows {
		incidents = append(incidents, models.Incident{
			Day:              int64(day),
			Month:            int64(month),
			Year:             row.Values[0].GetInt(),
			Summary:          row.Values[1].GetString_(),
			IncidentInDetail: row.Values[2].GetString_(),
			IncidentType:     models.Event,
		})
	}
	return incidents, nil
}

func (c *CassandraClient) FetchBirthsOnThisDay(day uint, month uint) ([]models.Incident, error) {
	stargateClient := clientPool.Get().(*client.StargateClient)
	defer clientPool.Put(stargateClient)
	fetchIncidentQuery := &datastax.Query{
		Cql: "SELECT year, summary, incident_in_summary FROM lear.births WHERE day = ? AND month = ?",
		Values: &datastax.Values{
			Values: []*datastax.Value{
				{
					Inner: &datastax.Value_Int{Int: int64(day)},
				},
				{
					Inner: &datastax.Value_Int{Int: int64(month)},
				},
			},
		},
	}

	resp, err := stargateClient.ExecuteQuery(fetchIncidentQuery)
	if err != nil {
		fmt.Printf("failed to fetch results from remote cassandra instance, reason %s\n", err.Error())
	}

	incidents := []models.Incident{}

	resultSet := resp.GetResultSet()

	for _, row := range resultSet.Rows {
		incidents = append(incidents, models.Incident{
			Day:              int64(day),
			Month:            int64(month),
			Year:             row.Values[0].GetInt(),
			Summary:          row.Values[1].GetString_(),
			IncidentInDetail: row.Values[2].GetString_(),
			IncidentType:     models.Birth,
		})
	}
	return incidents, nil
}

func (c *CassandraClient) FetchDeathsOnThisDay(day uint, month uint) ([]models.Incident, error) {
	stargateClient := clientPool.Get().(*client.StargateClient)
	defer clientPool.Put(stargateClient)
	fetchIncidentQuery := &datastax.Query{
		Cql: "SELECT year, summary, incident_in_summary FROM lear.deaths WHERE day = ? AND month = ?",
		Values: &datastax.Values{
			Values: []*datastax.Value{
				{
					Inner: &datastax.Value_Int{Int: int64(day)},
				},
				{
					Inner: &datastax.Value_Int{Int: int64(month)},
				},
			},
		},
	}

	resp, err := stargateClient.ExecuteQuery(fetchIncidentQuery)
	if err != nil {
		fmt.Printf("failed to fetch results from remote cassandra instance, reason %s\n", err.Error())
	}

	incidents := []models.Incident{}

	resultSet := resp.GetResultSet()

	for _, row := range resultSet.Rows {
		incidents = append(incidents, models.Incident{
			Day:              int64(day),
			Month:            int64(month),
			Year:             row.Values[0].GetInt(),
			Summary:          row.Values[1].GetString_(),
			IncidentInDetail: row.Values[2].GetString_(),
			IncidentType:     models.Death,
		})
	}
	return incidents, nil
}

func (c *CassandraClient) FetchHolidaysOnThisDay(day uint, month uint) ([]models.Incident, error) {
	stargateClient := clientPool.Get().(*client.StargateClient)
	defer clientPool.Put(stargateClient)
	fetchIncidentQuery := &datastax.Query{
		Cql: "SELECT year, summary, incident_in_summary FROM lear.holidays WHERE day = ? AND month = ?",
		Values: &datastax.Values{
			Values: []*datastax.Value{
				{
					Inner: &datastax.Value_Int{Int: int64(day)},
				},
				{
					Inner: &datastax.Value_Int{Int: int64(month)},
				},
			},
		},
	}

	resp, err := stargateClient.ExecuteQuery(fetchIncidentQuery)
	if err != nil {
		fmt.Printf("failed to fetch results from remote cassandra instance, reason %s\n", err.Error())
	}

	incidents := []models.Incident{}

	resultSet := resp.GetResultSet()

	for _, row := range resultSet.Rows {
		incidents = append(incidents, models.Incident{
			Day:              int64(day),
			Month:            int64(month),
			Year:             row.Values[0].GetInt(),
			Summary:          row.Values[1].GetString_(),
			IncidentInDetail: row.Values[2].GetString_(),
			IncidentType:     models.Holiday,
		})
	}
	return incidents, nil
}

func (c *CassandraClient) AddIncident(incident models.Incident) (models.Incident, error) {

	var stargate *client.StargateClient

	func() {
		stargate = clientPool.Get().(*client.StargateClient)
		atomic.AddInt64(&allocatedConnectionsToCassandra, 1)
	}()

	defer func() {
		clientPool.Put(stargate)
		atomic.AddInt64(&allocatedConnectionsToCassandra, -1)
	}()

	// switch {
	// case int64(time.Now().Year()) < incident.Year:
	// 	return models.Incident{}, fmt.Errorf("the event cannot be in the future, the year value is %v", incident.Month)
	// case incident.Month <= 0 || incident.Month > 12:
	// 	return models.Incident{}, fmt.Errorf("the month cannot be less than 1 as it has the value %v", incident.Month)
	// case incident.Day <= 0 || incident.Day > 31:
	// 	return models.Incident{}, fmt.Errorf("the day cannot be less than 1 as it has the value %v", incident.Day)
	// case incident.Summary == "":
	// 	return models.Incident{}, errors.New("the incident summary cannot be empty")
	// case incident.IncidentInDetail == "":
	// 	return models.Incident{}, errors.New("the incident in detail cannot be empty")
	// case incident.IncidentType != models.Event || incident.IncidentType != models.Birth || incident.IncidentType != models.Death || incident.IncidentType != models.Holiday:
	// 	return models.Incident{}, errors.New("unknown incident type")
	// }

	CqlInsertStatement := fmt.Sprintf(`INSERT INTO lear.%ss (day, month, year , summary , incident_in_detail ) VALUES ( ?, ?, ?, ?, ? );`, &incident.IncidentType)

	_, err := stargate.ExecuteQueryWithContext(&datastax.Query{
		Cql: CqlInsertStatement,
		Values: &datastax.Values{Values: []*datastax.Value{
			{Inner: &datastax.Value_Int{Int: int64(incident.Day)}},
			{Inner: &datastax.Value_Int{Int: int64(incident.Month)}},
			{Inner: &datastax.Value_Int{Int: int64(incident.Year)}},
			{Inner: &datastax.Value_String_{String_: incident.Summary}},
			{Inner: &datastax.Value_String_{String_: incident.IncidentInDetail}},
		}},
	}, context.Background())
	if err != nil {
		log.Fatalf("error, unable to insert incident into its table, reason: %v\n", err.Error())
	}

	return incident, nil
}

func (c *CassandraClient) BulkInsertIncidents(incidents models.Incident) ([]models.Incident, error) {
	//TODO implement me
	panic("implement me")
}
