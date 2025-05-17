package postgres

import (
	"database/sql"
	"time"

	"github.com/abiopereira/sw-criciuma/internal/domain"
	_ "github.com/lib/pq"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) domain.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) Create(message *domain.Message) error {
	query := `
		INSERT INTO messages (id, type, content, duration, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, message.ID, message.Type, message.Content, message.Duration, message.CreatedAt)
	return err
}

func (r *messageRepository) GetAll() ([]domain.Message, error) {
	query := `
		SELECT id, type, content, duration, created_at
		FROM messages
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		var createdAt time.Time
		err := rows.Scan(&msg.ID, &msg.Type, &msg.Content, &msg.Duration, &createdAt)
		if err != nil {
			return nil, err
		}
		msg.CreatedAt = createdAt
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *messageRepository) Delete(id string) error {
	query := `
		DELETE FROM messages
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
