package repository

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"database/sql"

	"github.com/google/uuid"
)

type ChatRepository struct {
	Db *configs.DatabaseConfig
}

func NewChatRepository(db *configs.DatabaseConfig) *ChatRepository {
	return &ChatRepository{Db: db}
}

func (chatRepository *ChatRepository) GetChatsRepository(senderID string, receiverID string) (*[]entity.ChatEntity, error) {
	query := `
		SELECT id, room_id, sender_id, receiver_id, message, created_at, read_at
		FROM chat
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $3 AND receiver_id = $4) 
		ORDER BY created_at ASC
	`

	rows, err := chatRepository.Db.DB.Connection.Query(query, senderID, receiverID, receiverID, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var readAt sql.NullTime
	var chats []entity.ChatEntity
	for rows.Next() {
		var job entity.ChatEntity
		err := rows.Scan(
			&job.ID,
			&job.RoomID,
			&job.SenderID,
			&job.ReceiverID,
			&job.Message,
			&job.CreatedAt,
			&readAt,
		)
		if err != nil {
			return nil, err
		}

		if readAt.Valid {
			job.ReadAt = readAt.Time
		}

		chats = append(chats, job)
	}

	return &chats, nil
}

func (chatRepository *ChatRepository) GetOrCreateRoomRepository(senderID, receiverID string) (string, error) {
	query := `
		SELECT id from rooms
		WHERE (sender_id = $1 AND receiver_id = $2 ) OR (sender_id = $3 AND receiver_id = $4);
	`
	var id string

	err := chatRepository.Db.DB.Connection.QueryRow(
		query,
		senderID,
		receiverID,
		receiverID,
		senderID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		id = uuid.NewString()
		query = `INSERT INTO rooms (id, sender_id, receiver_id) VALUES ($1, $2, $3)`
		_, err := chatRepository.Db.DB.Connection.Exec(
			query,
			senderID,
			receiverID,
		)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return id, nil
}

func (chatRepository *ChatRepository) CreateChatRepository(chat *entity.ChatEntity) (*entity.ChatEntity, error) {
	query := `
		INSERT INTO chat (id, room_id, sender_id, receiver_id, message)
		VALUES ($1, $2, $3, $4);
	`

	err := chatRepository.Db.DB.Connection.QueryRow(
		query,
		chat.ID,
		chat.RoomID,
		chat.SenderID,
		chat.ReceiverID,
		chat.Message,
	).Scan(
		&chat.ID,
		&chat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
