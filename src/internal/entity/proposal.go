package entity

import "time"

type ProposalEntity struct {
	ID            int64     `db:"id" json:"id"`
	JobID         int64     `db:"job_id" json:"job_id"`
	FreelancerID  int64     `db:"freelancer_id" json:"freelancer_id"`
	ProposalText  *string   `db:"proposal_text" json:"proposal_text,omitempty"`
	ProposedPrice *float64  `db:"proposed_price" json:"proposed_price,omitempty"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
