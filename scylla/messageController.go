package scylla

import (
	"org_chart/initializers"
	"org_chart/models"
)

func InsertMessage(msg *models.Message) error {
	query := `INSERT INTO messages (id, conversation_key, sender_id, receiver_id, content, sent_at, read_at) VALUES (?, ?, ?, ?, ?, ?, ?);`
	return initializers.Session.Query(query,
		msg.Id,
		msg.ConversationKey,
		msg.SenderId,
		msg.ReceiverId,
		msg.Content,
		msg.SentAt,
		nil,
	).Exec()
}

func FetchMessages(string) {
	// query := `SELECT * FROM messages `
	// initializers.Session.Query(query)
}
