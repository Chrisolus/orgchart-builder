package scylla

import (
	"errors"
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

func FetchMessages(conv_key string) ([]models.Message, error) {
	query := `SELECT id, conversation_key, sender_id, receiver_id, content, sent_at, read_at 
              FROM messages 
              WHERE conversation_key = ?;`

	iter := initializers.Session.Query(query, conv_key).Iter()

	var messages []models.Message
	var msg models.Message

	for iter.Scan(
		&msg.Id,
		&msg.ConversationKey,
		&msg.SenderId,
		&msg.ReceiverId,
		&msg.Content,
		&msg.SentAt,
		&msg.ReadAt,
	) {
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	if messages == nil {
		return nil, errors.New("no conversation exists")
	}
	return messages, nil
}
