package initializers

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitScyllaConnection(scylla_dsn string) {
	address, keyspace := filepath.Dir(scylla_dsn), filepath.Base(scylla_dsn)

	cluster := gocql.NewCluster(address)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = keyspace

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Scylla Connection error: ", err.Error())
	}

	initTable(keyspace)
}

func initTable(keyspace string) {
	query := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};", keyspace)
	err := Session.Query(query).Exec()
	if err != nil {
		log.Fatal("Scylla Keyspace Error: ", err.Error())
	} else {
		log.Printf("Keyspace: '%v' Created\n", keyspace)
	}

	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v.connections (
							id uuid, 
							client_id int, 
							is_active boolean, 
							connected_at timestamp, 
							disconnected_at timestamp, 
							PRIMARY KEY (client_id, id) );`,
		keyspace)
	err = Session.Query(query).Exec()
	if err != nil {
		log.Fatal("Connection Table Creation Error: ", err.Error())
	} else {
		log.Println("Connections Table Created")
	}

	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v.messages ( 
							id uuid, 
							conversation_key text, 
							sender_id int, 
							receiver_id int, 
							content text, 
							sent_at timestamp, 
							read_at timestamp, 
							PRIMARY KEY ((conversation_key), sent_at, id)) WITH CLUSTERING ORDER BY (sent_at ASC);`,
		keyspace)
	err = Session.Query(query).Exec()
	if err != nil {
		log.Fatal("Messages Table Creation Error: ", err.Error())
	} else {
		log.Println("Messages Table created")
	}
}

func Close() {
	Session.Close()
}
