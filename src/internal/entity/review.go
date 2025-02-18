package entity

import "time"

type ReviewEntity struct {
	ID         int64     `db:"id" json:"id"`
	JobID      int64     `db:"job_id" json:"job_id"`
	ReviewerID int64     `db:"reviewer_id" json:"reviewer_id"`
	Rating     int       `db:"rating" json:"rating"`
	Comment    *string   `db:"comment" json:"comment,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
