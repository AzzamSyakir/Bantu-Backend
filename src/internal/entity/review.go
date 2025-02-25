package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReviewEntity struct {
	ID         uuid.UUID `db:"id" json:"id"`
	JobID      uuid.UUID `db:"job_id" json:"job_id"`
	ReviewerID uuid.UUID `db:"reviewer_id" json:"reviewer_id"`
	Rating     int       `db:"rating" json:"rating"`
	Comment    string    `db:"comment" json:"comment,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
