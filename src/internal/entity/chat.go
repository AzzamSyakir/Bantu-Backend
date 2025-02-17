package entity

import "time"

type ChatEntity struct {
	ID         int64      `db:"id" json:"id"`
	SenderID   int64      `db:"sender_id" json:"sender_id"`
	ReceiverID int64      `db:"receiver_id" json:"receiver_id"`
	Message    *string    `db:"message" json:"message,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ReadAt     *time.Time `db:"read_at" json:"read_at,omitempty"`
}
