package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"org_chart/models"
	"org_chart/psql"
	"org_chart/scylla"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID      uint
	ConnID  gocql.UUID
	Conn    *websocket.Conn
	Receive chan []byte
	Hub     *Hub
}

type Hub struct {
	Clients map[uint]map[gocql.UUID]*Client
	Forward chan []byte
	Join    chan *Client
	Leave   chan *Client
}

var H *Hub

func InitHub() {
	H = &Hub{
		Clients: make(map[uint]map[gocql.UUID]*Client),
		Forward: make(chan []byte),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
	}
	H.Run()
}

func (hub *Hub) Run() {
	fmt.Println("HUB is active")
	for {
		select {
		case client := <-hub.Join:
			if hub.Clients[client.ID] == nil {
				hub.Clients[client.ID] = make(map[gocql.UUID]*Client)
			}
			hub.Clients[client.ID][client.ConnID] = client
			if err := scylla.AddConnection(client.ID, client.ConnID); err != nil {
				log.Println("connection activation error at db level: ", err.Error())
			}
		case client := <-hub.Leave:
			if conns, ok := hub.Clients[client.ID]; ok {
				delete(conns, client.ConnID)
				if len(conns) == 0 {
					delete(hub.Clients, client.ID)
				}
				close(client.Receive)
				client.Conn.Close()
				if err := scylla.RemoveConnection(client.ID, client.ConnID); err != nil {
					log.Println("deactivation error at db level: ", err.Error())
				}
			}
		case message := <-hub.Forward:
			var byteMsg []byte
			_, receiver, byteMsg := BuildMessageForDispatch(message, hub)
			for _, con := range receiver {
				if scylla.IsActive(con.ID) {
					con.Receive <- byteMsg
				}
			}
		}
	}
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Leave <- c
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		if err = validateIncomingMessage(message); err != nil {
			log.Println("Validation error: ", err.Error())
			response, _ := json.Marshal(models.Error{Type: "error", Content: "Invalid message structure"})
			c.Receive <- response
		} else {
			c.Hub.Forward <- message
		}
	}
}

func (c *Client) Write() {
	for message := range c.Receive {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}

}

func validateIncomingMessage(message []byte) error {
	var msg models.IncomingMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return errors.New("cannot parse message")
	}
	if msg.SenderId == 0 || msg.ReceiverId == 0 {
		log.Println(msg)
		return errors.New("sender_id / receiver_id fields cannot be empty")
	}

	if !psql.IsValidUserId(msg.SenderId) || !psql.IsValidUserId(msg.ReceiverId) {
		return errors.New("sender / receiver doesn't exist")
	}
	return nil
}

func BuildMessageForDispatch(msg []byte, hub *Hub) (map[gocql.UUID]*Client, map[gocql.UUID]*Client, []byte) {
	var incoming models.IncomingMessage
	json.Unmarshal(msg, &incoming)
	message := models.Message{
		Id:              gocql.TimeUUID(),
		Type:            "message",
		ConversationKey: GetConvKey(incoming.SenderId, incoming.ReceiverId),
		SenderId:        incoming.SenderId,
		ReceiverId:      incoming.ReceiverId,
		Content:         incoming.Content,
		SentAt:          time.Now(),
	}
	msgByte, _ := json.Marshal(message)

	if err := scylla.InsertMessage(&message); err != nil {
		log.Println("Message insertion failed: ", err.Error())
	}
	return hub.Clients[message.SenderId], hub.Clients[message.ReceiverId], msgByte
}

func GetConvKey(id1, id2 uint) string {
	return fmt.Sprintf("%d_%d", min(id1, id2), max(id1, id2))
}
