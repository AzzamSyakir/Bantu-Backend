package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProposalEntity struct {
	ID            uuid.UUID `db:"id" json:"id"`
	JobID         uuid.UUID `db:"job_id" json:"job_id"`
	FreelancerID  uuid.UUID `db:"freelancer_id" json:"freelancer_id"`
	ProposalText  *string   `db:"proposal_text" json:"proposal_text,omitempty"`
	ProposedPrice *float64  `db:"proposed_price" json:"proposed_price,omitempty"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
