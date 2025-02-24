package entity

import (
	"time"
)

type ChatEntity struct {
	ID         string     `db:"id" json:"id"`
	SenderID   string     `db:"sender_id" json:"sender_id"`
	ReceiverID string     `db:"receiver_id" json:"receiver_id"`
	Message    *string    `db:"message" json:"message,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ReadAt     *time.Time `db:"read_at" json:"read_at,omitempty"`
}
