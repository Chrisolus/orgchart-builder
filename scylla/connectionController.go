package scylla

import (
	"log"
	"org_chart/initializers"
	"time"

	"github.com/gocql/gocql"
)

type Connection struct {
	Id             gocql.UUID `json:"id"`
	ClientId       uint       `json:"client_id"`
	IsActive       bool       `json:"is_active"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt time.Time  `json:"disconnected_at"`
}

func AddConnection(clientId uint, connId gocql.UUID) error {
	con := Connection{
		Id:       connId,
		ClientId: clientId,
	}

	query := "INSERT INTO connections (id, client_id, is_active, connected_at) VALUES (?,?,?,?);"
	return initializers.Session.Query(query, con.Id, con.ClientId, true, time.Now()).Exec()

}

func RemoveConnection(clientId uint, connId gocql.UUID) error {
	query := "UPDATE connections SET is_active = false, disconnected_at = ? WHERE client_id = ? AND id = ?;"
	return initializers.Session.Query(query, time.Now(), clientId, connId).Exec()
}

func IsActive(clientId uint) bool {
	var connID gocql.UUID

	query := `SELECT id FROM connections 
          WHERE client_id = ? AND is_active = true 
          ALLOW FILTERING;`

	err := initializers.Session.Query(query, clientId).Consistency(gocql.One).Scan(&connID)

	if err == gocql.ErrNotFound {
		return false
	}
	if err != nil {
		log.Println("isactive error: ", err.Error())
		return false
	}

	return true
}
