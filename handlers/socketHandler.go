package handlers

import (
	"log"
	"net/http"
	"org_chart/hub"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "upgrade failed"})
		return
	}
	log.Println("Socket open for:")
	client := &hub.Client{ID: c.GetUint("user_id"), ConnID: gocql.TimeUUID(), Conn: conn, Hub: hub.H, Receive: make(chan []byte)}
	hub.H.Join <- client
	log.Printf("Client %v joined \n", client.ID)

	go client.Write()
	go client.Read()
}
