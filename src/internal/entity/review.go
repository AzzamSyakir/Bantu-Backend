package entity

import (
	"time"
)

type ReviewEntity struct {
	ID         string    `db:"id" json:"id"`
	JobID      string    `db:"job_id" json:"job_id"`
	ReviewerID string    `db:"reviewer_id" json:"reviewer_id"`
	Rating     int       `db:"rating" json:"rating"`
	Comment    *string   `db:"comment" json:"comment,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
