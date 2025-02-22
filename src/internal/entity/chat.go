package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatEntity struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	SenderID   uuid.UUID  `db:"sender_id" json:"sender_id"`
	ReceiverID uuid.UUID  `db:"receiver_id" json:"receiver_id"`
	Message    *string    `db:"message" json:"message,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ReadAt     *time.Time `db:"read_at" json:"read_at,omitempty"`
}
